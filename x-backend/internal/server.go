package internal

import (
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/cache"
	"github.com/jourloy/X-Backend/internal/handlers"
	"github.com/jourloy/X-Backend/internal/middlewares"
	"github.com/jourloy/X-Backend/internal/storage"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[gin]`,
		Level:  log.DebugLevel,
	})
)

func StartServer() {
	totalTime := time.Now()
	tempTime := time.Now()

	// Выключаем логгер
	gin.DefaultWriter = io.Discard

	// Инициализация хранилища
	storage.InitDB()
	logger.Debug(`Storage initialized`, `latency`, time.Since(tempTime))
	tempTime = time.Now()

	// Инициализация кэша
	cache.InitCache()
	logger.Debug(`Cache initialized`, `latency`, time.Since(tempTime))
	tempTime = time.Now()

	r := gin.New()

	// Middlewares
	r.Use(middlewares.Logger())
	r.Use(gin.Recovery())
	r.Use(middlewares.API())

	// Инициализация хендлеров
	handlers.Init(r)
	logger.Debug(`Handlers initialized`, `latency`, time.Since(tempTime))

	// Запуск сервера
	logger.Info(`Server started`, `port`, 3001, `latency (total)`, time.Since(totalTime))
	if err := r.Run(`0.0.0.0:3001`); err != nil {
		logger.Fatal(err)
	}
}
