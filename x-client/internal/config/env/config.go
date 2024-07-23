package env

import (
	"os"
)

var API string

func ParseENV() {
	if env, exist := os.LookupEnv(`API_KEY`); exist {
		API = env
	}
}
