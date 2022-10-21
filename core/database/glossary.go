package database

import (
	"time"

	"github.com/byvko-dev/am-core/mongodb/driver"

	"github.com/byvko-dev/am-types/stats/v3"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

const collectionVehiclesGlossary = "vehicle-glossary"
const collectionAchievementsGlossary = "achievement-glossary"

func GetVehiclesInfo(ids ...int) ([]stats.VehicleInfo, error) {
	client, err := driver.NewClient()
	if err != nil {
		return nil, err
	}

	filter := bson.M{}
	filter["tank_id"] = bson.M{"$in": ids}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data []stats.VehicleInfo
	cur, err := client.Raw(collectionVehiclesGlossary).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	return data, cur.All(ctx, &data)
}

func GetAchievementInfo(achievementId int) (stats.AchievementInfo, error) {
	client, err := driver.NewClient()
	if err != nil {
		return stats.AchievementInfo{}, err
	}

	var data stats.AchievementInfo
	filter := make(map[string]interface{})
	filter["achievement_id"] = achievementId
	return data, client.GetDocumentWithFilter(collectionAchievementsGlossary, filter, &data)
}
