package config

import "github.com/byvko-dev/am-core/helpers/env"

var ProxyHost = env.MustGetString("WG_PROXY_HOST")
