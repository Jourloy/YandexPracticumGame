package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/guards"
	"github.com/jourloy/X-Backend/internal/modules/sector"
)

func InitSector(r *gin.Engine) {
	h := initSectorHandlers()

	logger.Debug(`┏ Sector`)

	gMain := r.Group(`sector`)
	h.mainRoutes(gMain)
}

type sectorHandlers struct{}

func initSectorHandlers() *sectorHandlers {
	return &sectorHandlers{}
}

func (h sectorHandlers) mainRoutes(g *gin.RouterGroup) {
	controller := sector.InitController()

	g.Use(guards.CheckAPI())
	g.Use(guards.CheckAdmin())

	g.POST(``, controller.Create)
	logger.Debug(`┣━ Add`, `method`, `POST`, `route`, g.BasePath()+`/`)

	g.GET(`all`, controller.GetAll)
	logger.Debug(`┗━ Add`, `method`, `GET`, `route`, g.BasePath()+`/`)
}
