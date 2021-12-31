package wargaming

import (
	"fmt"
	"os"
)

type wargamingEndpoint string

func (e wargamingEndpoint) Fmt(args ...interface{}) string {
	return fmt.Sprintf(string(e), args...)
}
func (e wargamingEndpoint) AddBaseURI(realm string) wargamingEndpoint {
	return wargamingEndpoint(os.Getenv("WARGAMING_PROXY_URI") + "/" + realm + string(e))
}

// Accounts

var accountInfoEndpoint wargamingEndpoint = "/accounts/%v"                      // Call .Fmt(accountID) to get the endpoint
var accountClanInfoEndpoint wargamingEndpoint = "/accounts/%v/clan"             // Call .Fmt(accountID) to get the endpoint
var accountVehiclesEndpoint wargamingEndpoint = "/accounts/%v/vehicles"         // Call .Fmt(accountID) to get the endpoint
var accountAchievementsEndpoint wargamingEndpoint = "/accounts/%v/achievements" // Call .Fmt(accountID) to get the endpoint

// Clans

var clansNameSearch wargamingEndpoint = "/clans/search?q=%v" // Call .Fmt(keyword) to get the endpoint
var clansIDSearch wargamingEndpoint = "/clans/%v"            // Call .Fmt(clanID) to get the endpoint
