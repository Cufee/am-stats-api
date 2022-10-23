package database

import (
	"context"
	"time"

	"github.com/byvko-dev/am-core/mongodb/driver"
	"github.com/byvko-dev/am-types/stats/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionAccounts = "accounts"

func FindAccountByID(id int) (stats.AccountInfo, error) {
	client, err := driver.NewClient()
	if err != nil {
		return stats.AccountInfo{}, err
	}
	filter := make(map[string]interface{})
	filter["account_id"] = id

	var account stats.AccountInfo
	return account, client.GetDocumentWithFilter(collectionAccounts, filter, &account)
}

func UpdateAccount(profile stats.AccountInfo) error {
	client, err := driver.NewClient()
	if err != nil {
		return err
	}
	filter := bson.M{}
	filter["account_id"] = profile.AccountID

	update := bson.M{}
	update["$set"] = profile

	opts := options.Update().SetUpsert(true)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.Raw(collectionAccounts).UpdateOne(ctx, filter, update, opts)
	return err
}
