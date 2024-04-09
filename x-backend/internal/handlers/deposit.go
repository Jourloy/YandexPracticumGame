package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/guards"
	"github.com/jourloy/X-Backend/internal/modules/deposit"
)

func InitDeposit(r *gin.Engine) {
	h := initDepositHandlers()

	logger.Debug(`┏ Deposit`)

	gMain := r.Group(`deposit`)
	h.mainRoutes(gMain)
}

type depositHandlers struct{}

func initDepositHandlers() *depositHandlers {
	return &depositHandlers{}
}

func (h depositHandlers) mainRoutes(g *gin.RouterGroup) {
	controller := deposit.Init()

	g.Use(guards.CheckAPI())

	g.POST(``, controller.Create)
	logger.Debug(`┗━ Add`, `method`, `POST`, `route`, g.BasePath()+`/`)
}
