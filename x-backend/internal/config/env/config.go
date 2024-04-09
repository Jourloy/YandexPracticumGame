package env

import (
	"os"
)

var DatabaseDSN string
var CacheDSN string

func ParseENV() {
	if env, exist := os.LookupEnv(`DATABASE_DSN`); exist {
		DatabaseDSN = env
	}

	if env, exist := os.LookupEnv(`CACHE_DSN`); exist {
		CacheDSN = env
	}
}
