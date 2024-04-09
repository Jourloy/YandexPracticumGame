package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/guards"
	"github.com/jourloy/X-Backend/internal/modules/account"
)

func InitAccount(r *gin.Engine) {
	h := initAccountHandlers()

	logger.Debug(`┏ Account`)

	gMain := r.Group(`account`)
	h.mainRoutes(gMain)
}

type accountHandlers struct{}

func initAccountHandlers() *accountHandlers {
	return &accountHandlers{}
}

func (h accountHandlers) mainRoutes(g *gin.RouterGroup) {
	controller := account.InitController()

	g.POST(``, controller.Create)
	logger.Debug(`┣━ Add`, `method`, `POST`, `route`, g.BasePath()+`/`)

	g.Use(guards.CheckAPI())

	g.GET(``, controller.GetMe)
	logger.Debug(`┗━ Add`, `method`, `GET`, `route`, g.BasePath()+`/`)
}
