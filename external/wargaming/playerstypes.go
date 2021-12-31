package wargaming

import (
	wgtypes "aftermath.link/repo/am-types/wargaming"
)

// Response from accountInfoEndpoint
type playersDataEndpointJSON struct {
	Data   map[string]wgtypes.StatsFrame `json:"data"`
	Status string                        `json:"status"`
	Error  struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
		Value   string `json:"value"`
	} `json:"error"`
}
