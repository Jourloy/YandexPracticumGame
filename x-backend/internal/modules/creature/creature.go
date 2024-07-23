package creature

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/config/errs"
	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/tools"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[creature]`,
		Level:  log.DebugLevel,
	})
)

type controller struct {
	service Service
}

// Init создает контроллер существа
func Init() *controller {
	return &controller{
		service: *InitService(),
	}
}

// Create создает существо
// Для admin
func (s *controller) Create(c *gin.Context) {
	accountID := c.GetString(`accountID`)

	// Парсинг body
	var body repositories.CreatureCreate
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

	// TODO
	c.JSON(200, gin.H{})
}

type spawnResponseSuccess struct {
	Creature repositories.Creature `json:"creature"`
}

type spawnResponseError struct {
	Error string `json:"error"`
}

// Spawn создает существо
func (s *controller) Spawn(c *gin.Context) {
	accountID := c.GetString(`accountID`)

	// Парсинг body
	var body spawnCreature
	if err := tools.ParseBody(c, &body); err != nil {
		logger.Error(err)
		c.JSON(
			errs.Errors.ParseError.Code,
			spawnResponseError{Error: errs.Errors.ParseError.Error.Error()},
		)
	}

	body.AccountID = accountID

	resp := s.service.Spawn(body)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(400, gin.H{`error`: resp.Err.Error()})
	}

	c.JSON(200, spawnResponseSuccess{Creature: *resp.Creature})
}

type moveResponseSuccess struct {
	Creature repositories.Creature `json:"creature"`
}

type moveResponseError struct {
	Error string `json:"error"`
}

// Move двигает существо
func (s *controller) Move(c *gin.Context) {
	accountID := c.GetString(`accountID`)

	// Парсинг body
	var body moveCreature
	if err := tools.ParseBody(c, &body); err != nil {
		logger.Error(err)
		c.JSON(
			errs.Errors.ParseError.Code,
			spawnResponseError{Error: errs.Errors.ParseError.Error.Error()},
		)
	}

	body.AccountID = accountID

	resp := s.service.Move(body)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(400, gin.H{`error`: resp.Err.Error()})
	}

	c.JSON(200, moveResponseSuccess{Creature: *resp.Creature})
}
