package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/guards"
	"github.com/jourloy/X-Backend/internal/modules/creature"
)

func InitCreature(r *gin.Engine) {
	h := initCreatureHandlers()

	logger.Debug(`┏ Creature`)

	gMain := r.Group(`creature`)
	h.mainRoutes(gMain)
}

type creatureHandlers struct{}

func initCreatureHandlers() *creatureHandlers {
	return &creatureHandlers{}
}

func (h creatureHandlers) mainRoutes(g *gin.RouterGroup) {
	controller := creature.Init()

	g.Use(guards.CheckAPI())

	g.POST(``, controller.Create)
	logger.Debug(`┣━ Add`, `method`, `POST`, `route`, g.BasePath()+`/`)

	g.POST(`spawn`, controller.Move)
	logger.Debug(`┣━ Add`, `method`, `POST`, `route`, g.BasePath()+`/spawn`)

	g.POST(`move`, controller.Move)
	logger.Debug(`┗━ Add`, `method`, `POST`, `route`, g.BasePath()+`/move`)
}
