package main

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"

	"github.com/jourloy/X-Backend/internal"
	"github.com/jourloy/X-Backend/internal/config/env"
)

// @Title X API
// @Description Документация для REST API игры
// @Version 1.0

// @BasePath /
// @Host localhost:3001

// @SecurityDefinitions.apikey ApiKeyAuth
// @In header
// @Name api-key

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear[`linux`] = func() {
		cmd := exec.Command(`clear`)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(`Your platform is unsupported`)
		}
	}
	clear[`darwin`] = func() {
		cmd := exec.Command(`clear`)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(`Your platform is unsupported`)
		}
	}
	clear[`windows`] = func() {
		cmd := exec.Command(`cmd`, `/c`, `cls`)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			panic(`Your platform is unsupported`)
		}
	}
}

func main() {
	// Очистка терминала
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic(`Your platform is unsupported`)
	}

	// Загрузка .env
	if err := godotenv.Load(`.env`); err != nil {
		log.Fatal(`Error loading .env file`)
	}

	// Парсинг .env
	env.ParseENV()

	// Старт сервера
	internal.StartServer()
}
