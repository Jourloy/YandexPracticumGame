package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/guards"
	"github.com/jourloy/X-Backend/internal/modules/building"
)

func InitBuilding(r *gin.Engine) {
	h := initBuildingHandlers()

	logger.Debug(`┏ Building`)

	gMain := r.Group(`building`)
	h.mainRoutes(gMain)
}

type buidingHandlers struct{}

func initBuildingHandlers() *buidingHandlers {
	return &buidingHandlers{}
}

func (h buidingHandlers) mainRoutes(g *gin.RouterGroup) {
	controller := building.InitController()

	g.Use(guards.CheckAPI())

	g.POST(``, controller.Create)
	logger.Debug(`┣━ Add`, `method`, `POST`, `route`, g.BasePath()+`/`)

	g.POST(`townhall`, controller.PlaceTownHall)
	logger.Debug(`┣━ Add`, `method`, `POST`, `route`, g.BasePath()+`/townhall`)

	g.GET(``, controller.GetOne)
	logger.Debug(`┣━ Add`, `method`, `GET`, `route`, g.BasePath()+`/`)

	g.GET(`all`, controller.GetAll)
	logger.Debug(`┗━ Add`, `method`, `GET`, `route`, g.BasePath()+`/all`)
}
