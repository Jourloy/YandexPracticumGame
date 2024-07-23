package handlers

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[handlers]`,
		Level:  log.DebugLevel,
	})
)

func Init(r *gin.Engine) {
	InitAccount(r)

	InitDeposit(r)
	InitResource(r)
	InitNode(r)
	InitSector(r)

	InitCreature(r)
	InitBuilding(r)
}
