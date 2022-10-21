package tests

import (
	"testing"

	"github.com/byvko-dev/am-stats-api/logic"
)

// Test GetPlayerSession
func TestLogicGetPlayerSession(t *testing.T) {
	accountId := 1013379500
	manual := false
	days := 0
	_, session, _, err := logic.GetPlayerSession(accountId, days, manual)
	if err != nil {
		t.Error(err)
	}
	t.Log("Total battles in session: ", session.Regular.Total.Battles+session.Rating.Total.Battles)
}
