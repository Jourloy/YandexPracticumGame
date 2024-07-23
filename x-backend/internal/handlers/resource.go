package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/guards"
	"github.com/jourloy/X-Backend/internal/modules/resource"
)

func InitResource(r *gin.Engine) {
	h := initResourceHandlers()

	logger.Debug(`┏ Resource`)

	gMain := r.Group(`resource`)
	h.mainRoutes(gMain)
}

type resourceHandlers struct{}

func initResourceHandlers() *resourceHandlers {
	return &resourceHandlers{}
}

func (h resourceHandlers) mainRoutes(g *gin.RouterGroup) {
	controller := resource.Init()

	g.Use(guards.CheckAPI())

	g.POST(``, controller.Create)
	logger.Debug(`┗━ Add`, `method`, `POST`, `route`, g.BasePath()+`/`)
}
