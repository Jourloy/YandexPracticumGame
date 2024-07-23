package operation

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/tools"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[operation]`,
		Level:  log.DebugLevel,
	})
)

type controller struct {
	service Service
}

// Init создает контроллер операции
func Init() *controller {
	return &controller{
		service: *InitService(),
	}
}

// Create создает операцию
func (s *controller) Create(c *gin.Context) {
	accountID := c.GetString(`accountID`)

	// Парсинг body
	var body repositories.OperationCreate
	if err := tools.ParseBody(c, &body); err != nil {
		logger.Error(`parse body error`)
		c.JSON(400, gin.H{`error`: `parse body error`})
		return
	}

	body.AccountID = accountID

	resp := s.service.Create(body)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(resp.Code, gin.H{`error`: resp.Err.Error()})
		return
	}

	c.JSON(resp.Code, gin.H{`error`: ``, `operation`: resp.Operation})
}
