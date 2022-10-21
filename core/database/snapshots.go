package database

import (
	"errors"
	"time"

	"github.com/byvko-dev/am-core/mongodb/driver"
	"github.com/byvko-dev/am-types/stats/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

const collectionSnapshots = "snapshots"

func GetPlayerSnapshot(accountId, days int, manual bool, currentBattles int) (stats.AccountSnapshot, error) {
	client, err := driver.NewClient()
	if err != nil {
		return stats.AccountSnapshot{}, err
	}

	var opts options.FindOneOptions
	opts.SetSort(bson.M{"created_at": -1})

	var snapshot stats.AccountSnapshot
	var filter bson.M = bson.M{}
	filter["is_manual"] = manual
	filter["account_id"] = accountId
	filter["total_battles"] = bson.M{"$lt": currentBattles}
	if days > 0 {
		// Find the oldest snapshot that was created before the specified number of days
		timestamp := time.Now().AddDate(0, 0, -days).Unix()
		filter["created_at"] = bson.M{"$gt": timestamp}
		opts.SetSort(bson.M{"created_at": 1})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return snapshot, client.Raw(collectionSnapshots).FindOne(ctx, filter).Decode(&snapshot)
}

func GetLastTotalBattles(accountId int, isManual bool) (int, error) {
	client, err := driver.NewClient()
	if err != nil {
		return 0, err
	}

	var target stats.AccountSnapshot
	filter := bson.M{"account_id": accountId, "is_manual": isManual}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var opts options.FindOneOptions
	opts.SetSort(bson.M{"created_at": -1})

	err = client.Raw(collectionSnapshots).FindOne(ctx, filter, &opts).Decode(&target)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return 0, nil
	}
	return target.TotalBattles, err
}
