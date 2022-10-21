package logic

import (
	"errors"
	"sync"
	"time"

	"github.com/byvko-dev/am-core/logs"
	"github.com/byvko-dev/am-core/stats/blitzstars/v1/types"
	"github.com/byvko-dev/am-stats-api/config"
	"github.com/byvko-dev/am-stats-api/core/database"
	"github.com/byvko-dev/am-stats-updates/calculations"
	e "github.com/byvko-dev/am-types/errors/v2"
	"github.com/byvko-dev/am-types/stats/v3"
	"github.com/byvko-dev/am-types/wargaming/v2/accounts"
	"github.com/byvko-dev/am-types/wargaming/v2/statistics"
	wg "github.com/cufee/am-wg-proxy-next/client"
	"go.mongodb.org/mongo-driver/mongo"
)

type data interface {
	stats.AccountSnapshot | statistics.AchievementsFrame | accounts.CompleteProfile | []statistics.VehicleStatsFrame | []types.TankAverages | []stats.VehicleInfo
}

type result[T data] struct {
	data T
	err  *e.Error
}

func GetPlayerSession(accountId, days int, manual bool) (profile stats.AccountInfo, sessionFrame stats.CompleteFrame, snapshotFrame stats.CompleteFrame, err *e.Error) {
	client := wg.NewClient(config.ProxyHost, time.Second*60)
	defer client.Close()

	start := time.Now()

	var wg sync.WaitGroup

	// Live player overall stats
	snapshotResult := make(chan result[stats.AccountSnapshot], 1)
	accountResult := make(chan result[accounts.CompleteProfile], 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		account, err := client.GetAccountByID(accountId)
		accountResult <- result[accounts.CompleteProfile]{data: account, err: err}

		// Player snapshot from DB (if exists)
		currentBattles := account.Statistics.All.Battles + account.Statistics.Rating.Battles
		snapshot, baseErr := database.GetPlayerSnapshot(accountId, days, manual, currentBattles)
		if baseErr != nil {
			if errors.Is(baseErr, mongo.ErrNoDocuments) {
				snapshotResult <- result[stats.AccountSnapshot]{data: stats.AccountSnapshot{}, err: nil}
				return
			}
			snapshotResult <- result[stats.AccountSnapshot]{err: e.Generic(baseErr, "failed to get player snapshot")}
			return
		}
		snapshotResult <- result[stats.AccountSnapshot]{data: snapshot}

	}()

	// Live player vehicle stats
	averagesResult := make(chan result[[]types.TankAverages], 1)
	glossaryResult := make(chan result[[]stats.VehicleInfo], 1)
	vehiclesResult := make(chan result[[]statistics.VehicleStatsFrame], 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		vehicles, err := client.GetAccountVehicles(accountId)
		vehiclesResult <- result[[]statistics.VehicleStatsFrame]{data: vehicles, err: err}

		var ids []int
		for _, vehicle := range vehicles {
			ids = append(ids, vehicle.TankID)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			averages, basicErr := database.GetTankAverages(ids...)
			if basicErr != nil {
				averagesResult <- result[[]types.TankAverages]{err: e.Generic(basicErr, "failed to get tank averages")}
				return
			}
			averagesResult <- result[[]types.TankAverages]{data: averages}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			info, basicErr := database.GetVehiclesInfo(ids...)
			if basicErr != nil {
				glossaryResult <- result[[]stats.VehicleInfo]{err: e.Generic(basicErr, "failed to get tank averages")}
				return
			}
			glossaryResult <- result[[]stats.VehicleInfo]{data: info}
		}()
	}()

	// Live player achievements
	achievementsResult := make(chan result[statistics.AchievementsFrame], 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		achievements, err := client.GetAccountAchievements(accountId)
		achievementsResult <- result[statistics.AchievementsFrame]{data: achievements, err: err}
	}()

	// TODO: Add support for vehicle achievements
	vehicleAchievements := make(map[int]statistics.AchievementsFrame)

	logs.Debug("Started routines: %v", time.Since(start))

	wg.Wait()
	close(snapshotResult)
	close(accountResult)
	close(vehiclesResult)
	close(achievementsResult)
	close(averagesResult)
	close(glossaryResult)

	snapshot := <-snapshotResult
	if snapshot.err != nil {
		err = snapshot.err
		return
	}

	account := <-accountResult
	if account.err != nil {
		err = account.err
		return
	}

	achievements := <-achievementsResult
	if achievements.err != nil {
		err = achievements.err
		return
	}

	vehicles := <-vehiclesResult
	if vehicles.err != nil {
		err = vehicles.err
		return
	}

	averagesSlice := <-averagesResult
	if averagesSlice.err != nil {
		err = averagesSlice.err
		return
	}
	averages := make(map[int]types.TankAverages)
	for _, average := range averagesSlice.data {
		averages[average.TankID] = average
	}

	glossarySlice := <-glossaryResult
	if glossarySlice.err != nil {
		err = glossarySlice.err
		return
	}
	glossaryData := make(map[int]stats.VehicleInfo)
	for _, data := range glossarySlice.data {
		glossaryData[data.TankID] = data
	}

	logs.Debug("Found %v tank vehicle info", len(glossaryData))

	logs.Debug("Received all replies: %v", time.Since(start))

	vehicleCutoffTime := snapshot.data.LastBattleTime // We don't need session for vehicles that were not played during this session
	if snapshot.data.TotalBattles == 0 {
		vehicleCutoffTime = int(time.Now().Add(time.Hour).Unix()) // If there is no snapshot, we need to get any vehicle stats
	}
	liveSnapshot, baseErr := calculations.AccountSnapshot(account.data, achievements.data, vehicles.data, vehicleAchievements, vehicleCutoffTime, averages, glossaryData)
	if baseErr != nil {
		err = e.Input(baseErr, "failed to calculate live player snapshot")
		return
	}

	logs.Debug("Calculated snapshot: %v", time.Since(start))

	if snapshot.data.TotalBattles == 0 {
		sessionFrame = stats.CompleteFrame{Regular: liveSnapshot.Stats.Regular, Rating: liveSnapshot.Stats.Rating}
		err = nil
		return
	}

	sessionRegular := liveSnapshot.Stats.Regular
	sessionRating := liveSnapshot.Stats.Rating

	sessionRegular.Subtract(&snapshot.data.Stats.Regular)
	sessionRating.Subtract(&snapshot.data.Stats.Rating)

	logs.Debug("Calculated session: %v", time.Since(start))

	sessionFrame = stats.CompleteFrame{
		Regular: sessionRegular,
		Rating:  sessionRating,
	}
	snapshotFrame = snapshot.data.Stats
	err = nil
	return
}
