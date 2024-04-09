package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/jourloy/x-client/internal/config/env"
	"github.com/jourloy/x-client/internal/view/input"
	"github.com/jourloy/x-client/internal/view/tabs"
)

var version string = `N/A`
var date string = `N/A`

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Level: log.DebugLevel,
	})
)

func main() {
	logger.Info(`Client info`, `version`, version, `build date`, date)

	// Загрузка .env
	if err := godotenv.Load(`.env`); err != nil {
		log.Error(`Error loading .env file`)
	} else {
		env.ParseENV()
	}

	if env.API == `` {
		input.Main()
	}

	tabs.Main()
}
