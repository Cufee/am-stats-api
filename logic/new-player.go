package logic

import (
	"time"

	"github.com/byvko-dev/am-core/logs"
	"github.com/byvko-dev/am-stats-api/config"
	"github.com/byvko-dev/am-stats-api/core/cache"
	"github.com/byvko-dev/am-stats-api/core/database"
	"github.com/byvko-dev/am-types/stats/v3"
	wg "github.com/cufee/am-wg-proxy-next/client"
	"github.com/cufee/am-wg-proxy-next/helpers"
)

func CheckPlayerCache(account stats.AccountInfo) {
	if account.AccountID == 0 {
		logs.Error("account id is 0")
		return
	}

	cached, _ := database.FindAccountByID(int(account.AccountID))
	if cached.AccountID == account.AccountID {
		return
	}

	client := wg.NewClient(config.ProxyHost, time.Second*10)
	defer client.Close()

	live, err := client.GetAccountByID(int(account.AccountID))
	if err != nil {
		logs.Error("Failed to get account from WG proxy: %v", err.Error())
		return
	}

	account.Nickname = live.Nickname
	account.AccountID = int(live.AccountID)
	account.Realm = helpers.RealmFromID(account.AccountID)

	clan, _ := client.GetAccountClan(int(account.AccountID))
	if clan.ClanID != 0 {
		account.Clan.JoinedAt = int(clan.JoinedAt)
		account.Clan.Name = clan.Clan.Name
		account.Clan.Role = clan.Role
		account.Clan.Tag = clan.Clan.Tag
		account.Clan.ID = int(clan.ClanID)
	}

	basicErr := database.UpdateAccount(account)
	if basicErr != nil {
		logs.Error("Failed to update account in db: %v", basicErr.Error())
		return
	}

	basicErr = cache.RecordPlayerSessions(account.Realm, false, account.AccountID)
	if basicErr != nil {
		logs.Error("Failed to record player sessions: %v", basicErr.Error())
		return
	}
}
