package database

import (
	"time"

	"github.com/byvko-dev/am-core/mongodb/driver"
	"github.com/byvko-dev/am-core/stats/blitzstars/v1/types"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
)

const collectionAverages = "vehicle-averages"

func GetTankAverages(ids ...int) ([]types.TankAverages, error) {
	client, err := driver.NewClient()
	if err != nil {
		return nil, err
	}

	filter := bson.M{}
	filter["tankId"] = bson.M{"$in": ids}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data []types.TankAverages
	cur, err := client.Raw(collectionAverages).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	return data, cur.All(ctx, &data)
}
