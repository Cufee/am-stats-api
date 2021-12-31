package wargaming

import (
	"fmt"
	"os"
	"testing"
)

func TestAddBaseURI(t *testing.T) {
	realm := "eu"
	player := 123321

	os.Setenv("WARGAMING_PROXY_URI", "http://wargaming.example.com")
	var e wargamingEndpoint = accountInfoEndpoint
	u := e.AddBaseURI(realm).Fmt(player)
	if u != fmt.Sprintf("http://wargaming.example.com/%v/accounts/%v/info", realm, player) {
		t.Errorf("wargamingEndpoint.AddBaseURI() failed: %v", u)
	}
}

func TestFmt(t *testing.T) {
	value := 12345
	var e wargamingEndpoint = "test/%v"

	result := e.Fmt(value)
	if result != "test/12345" {
		t.Errorf("wargamingEndpoint.Fmt() failed: %v", result)
	}
}
