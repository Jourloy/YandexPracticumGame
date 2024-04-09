package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/guards"
	"github.com/jourloy/X-Backend/internal/modules/node"
)

func InitNode(r *gin.Engine) {
	h := initNodeHandlers()

	gMain := r.Group(`node`)
	h.mainRoutes(gMain)
}

type nodeHandlers struct{}

func initNodeHandlers() *nodeHandlers {
	return &nodeHandlers{}
}

func (h nodeHandlers) mainRoutes(g *gin.RouterGroup) {
	node.InitController()

	g.Use(guards.CheckAPI())
}
