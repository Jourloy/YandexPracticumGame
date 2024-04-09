package account

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
		Prefix: `[account]`,
		Level:  log.DebugLevel,
	})
)

type controller struct{}

// InitController создает сервис аккаунта
func InitController() *controller {
	initService()

	return &controller{}
}

type createResponseSuccess struct {
	Account repositories.Account `json:"account"`
}

type createResponseError struct {
	Error string `json:"error"`
}

// Create создает аккаунт
func (s *controller) Create(c *gin.Context) {
	var body repositories.AccountCreate
	if err := tools.ParseBody(c, &body); err != nil {
		logger.Error(err)
		c.JSON(
			errs.Errors.ParseError.Code,
			createResponseError{Error: errs.Errors.ParseError.Error.Error()},
		)
		return
	}

	resp := Service.Create(body)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(
			resp.ErrResp.Code,
			createResponseError{Error: resp.ErrResp.Error.Error()},
		)
		return
	}

	c.JSON(200, createResponseSuccess{Account: *resp.Account})
}

type getResponseSuccess struct {
	Account repositories.Account `json:"account"`
}

type getResponseError struct {
	Error string `json:"error"`
}

// GetMe получает аккаунт авторизованного пользователя
func (s *controller) GetMe(c *gin.Context) {
	a, exist := c.Get(`account`)
	if !exist {
		logger.Error(errs.Errors.Unathorized)
		c.JSON(
			errs.Errors.Unathorized.Code,
			createResponseError{Error: errs.Errors.Unathorized.Error.Error()},
		)
		return
	}

	account, ok := a.(repositories.Account)
	if !ok {
		logger.Error(errs.Errors.InvalidAccount)
		c.JSON(
			errs.Errors.InvalidAccount.Code,
			createResponseError{Error: errs.Errors.InvalidAccount.Error.Error()},
		)
		return
	}

	resp := Service.GetOne(&repositories.AccountGet{ID: &account.ID})
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(
			resp.ErrResp.Code,
			createResponseError{Error: resp.ErrResp.Error.Error()},
		)
		return
	}

	c.JSON(200, getResponseSuccess{Account: *resp.Account})
}
