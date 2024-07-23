package storage

import (
	"os"

	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"

	"github.com/jourloy/X-Backend/internal/config/env"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[storage]`,
		Level:  log.DebugLevel,
	})
)

var Database *gorm.DB

// InitDB подключается к базе данных
func InitDB() {
	db, err := gorm.Open(postgres.Open(env.DatabaseDSN), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 lg.Default.LogMode(lg.Silent),
	})
	if err != nil {
		logger.Fatal(`Failed to connect database`)
	}
	Database = db
}
