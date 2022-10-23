package cache

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/byvko-dev/am-core/helpers/env"
	"github.com/byvko-dev/am-core/helpers/requests"
)

var apiURL = env.MustGetString("CACHE_API_URL")

func RecordPlayerSessions(realm string, manual bool, ids ...int) error {
	var idsStr []string
	for _, id := range ids {
		idsStr = append(idsStr, fmt.Sprint(id))
	}

	endpoint, err := url.Parse(fmt.Sprintf("%s/realm/%v/sessions/players", apiURL, realm))
	if err != nil {
		return err
	}

	q := endpoint.Query()
	q.Set("ids", strings.Join(idsStr, ","))
	q.Set("manual", fmt.Sprint(manual))
	endpoint.RawQuery = q.Encode()

	status, err := requests.Send(endpoint.String(), "GET", nil, nil, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("unexpected status code: %d", status)
	}
	return nil

}
