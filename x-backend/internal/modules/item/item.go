package item

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
		Prefix: `[item]`,
		Level:  log.DebugLevel,
	})
)

type controller struct {
	service Service
}

// Init создает сервис предмета
func Init() *controller {
	return &controller{
		service: *InitService(),
	}
}

type createResponseSuccess struct {
	Error string            `json:"error"`
	Item  repositories.Item `json:"item"`
}

type createResponseError struct {
	Error string `json:"error"`
}

// Create создает предмет
func (s *controller) Create(c *gin.Context) {
	a, exist := c.Get(`account`)
	if !exist {
		logger.Error(errs.Errors.Unathorized.Error)
		c.JSON(
			errs.Errors.Unathorized.Code,
			createResponseError{Error: errs.Errors.Unathorized.Error.Error()},
		)
		return
	}

	account, ok := a.(repositories.Account)
	if !ok {
		logger.Error(errs.Errors.InvalidAccount.Error)
		c.JSON(
			errs.Errors.InvalidAccount.Code,
			createResponseError{Error: errs.Errors.InvalidAccount.Error.Error()},
		)
		return
	}

	// Парсинг body
	var body repositories.ItemCreate
	if err := tools.ParseBody(c, &body); err != nil {
		logger.Error(errs.Errors.ParseError.Error)
		c.JSON(
			errs.Errors.ParseError.Code,
			createResponseError{Error: errs.Errors.ParseError.Error.Error()},
		)
		return
	}
	body.CreatorID = account.ID

	resp := s.service.Create(body)
	if resp.Err != nil {
		logger.Error(errs.Errors.ParseError.Error)
		c.JSON(
			resp.ErrResp.Code,
			createResponseError{Error: resp.ErrResp.Error.Error()},
		)
	}

	c.JSON(200, createResponseSuccess{Item: *resp.Item})
}
