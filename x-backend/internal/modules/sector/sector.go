package sector

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
		Prefix: `[sector]`,
		Level:  log.DebugLevel,
	})
)

type controller struct{}

func InitController() *controller {
	return &controller{}
}

func (s *controller) Create(c *gin.Context) {
	// Парсинг body
	var body CreateOptions
	if err := tools.ParseBody(c, &body); err != nil {
		logger.Error(`parse body error`)
		c.JSON(400, gin.H{`error`: `parse body error`})
		return
	}

	resp := Service.Create(body)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(400, gin.H{`error`: resp.Err.Error()})
		return
	}

	c.JSON(200, gin.H{`error`: ``})
}

type getAllResponseSuccess struct {
	Sectors []repositories.Sector `json:"sectors"`
	Count   int
}

type getAllResponseError struct {
	Error string `json:"error"`
}

func (s *controller) GetAll(c *gin.Context) {

	// Создание фильтров
	var query repositories.SectorGet
	if err := c.Bind(&query); err != nil {
		logger.Error(err)
		c.JSON(
			errs.Errors.ParseError.Code,
			getAllResponseError{Error: errs.Errors.ParseError.Error.Error()},
		)
		return
	}

	resp := Service.GetAll(query)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(
			resp.ErrResp.Code,
			getAllResponseError{Error: resp.ErrResp.Error.Error()},
		)
		return
	}

	c.JSON(200, getAllResponseSuccess{Sectors: *resp.Sectors, Count: len(*resp.Sectors)})
}
