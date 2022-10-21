package database

import (
	"github.com/byvko-dev/am-core/mongodb/driver"
	"github.com/byvko-dev/am-types/stats/v3"
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
