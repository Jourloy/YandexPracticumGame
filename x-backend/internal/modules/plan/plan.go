package plan

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/tools"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[plan]`,
		Level:  log.DebugLevel,
	})
)

type controller struct {
	service Service
}

// Init создает контроллер планируемой постройки
func Init() *controller {
	return &controller{
		service: *InitService(),
	}
}

// Create создает планируемую постройку
func (s *controller) Create(c *gin.Context) {
	accountID := c.GetString(`accountID`)

	// Парсинг body
	var body repositories.PlanCreate
	if err := tools.ParseBody(c, &body); err != nil {
		logger.Error(`Parse body error`)
		c.JSON(400, gin.H{`error`: `Parse body error`})
	}

	body.AccountID = accountID

	resp := s.service.Create(body)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(400, gin.H{`error`: resp.Err.Error()})
	}

	c.JSON(200, gin.H{`error`: ``})
}
