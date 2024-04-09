package node

import (
	"os"

	"github.com/charmbracelet/log"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[node]`,
		Level:  log.DebugLevel,
	})
)

type controller struct {
	service Service
}

// InitController создает сервис аккаунта
func InitController() *controller {
	return &controller{
		service: *InitService(),
	}
}
