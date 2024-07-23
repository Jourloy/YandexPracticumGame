package building

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
		Prefix: `[building]`,
		Level:  log.DebugLevel,
	})
)

type controller struct{}

// InitController создает контроллер постройки
func InitController() *controller {
	return &controller{}
}

type createResponseSuccess struct {
	Building repositories.Building `json:"building"`
}

type createResponseError struct {
	Error string `json:"error"`
}

func (s *controller) Create(c *gin.Context) {
	accountID := c.GetString(`accountID`)

	// Парсинг body
	var body repositories.BuildingCreate
	if err := tools.ParseBody(c, &body); err != nil {
		logger.Error(err)
		c.JSON(
			errs.Errors.ParseError.Code,
			createResponseError{Error: errs.Errors.ParseError.Error.Error()},
		)
		return
	}

	body.AccountID = accountID

	resp := Service.Create(body)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(
			resp.ErrResp.Code,
			createResponseError{Error: resp.ErrResp.Error.Error()},
		)
		return
	}

	c.JSON(200, createResponseSuccess{Building: *resp.Building})
}

type getResponseSuccess struct {
	Building repositories.Building `json:"building"`
}

type getResponseError struct {
	Error string `json:"error"`
}

func (s *controller) GetOne(c *gin.Context) {
	// Парсинг body
	var query repositories.BuildingGet
	if err := c.Bind(&query); err != nil {
		logger.Error(err)
		c.JSON(
			errs.Errors.ParseError.Code,
			getResponseError{Error: errs.Errors.ParseError.Error.Error()},
		)
		return
	}

	resp := Service.GetOne(&query)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(
			resp.ErrResp.Code,
			getResponseError{Error: resp.ErrResp.Error.Error()},
		)
		return
	}

	c.JSON(200, getResponseSuccess{Building: *resp.Building})
}

type getAllResponseSuccess struct {
	Buildings []repositories.Building `json:"buildings"`
	Count     int
}

type getAllResponseError struct {
	Error string `json:"error"`
}

// GetAll возвращает все постройки
func (s *controller) GetAll(c *gin.Context) {
	// Парсинг body
	var query repositories.BuildingGet
	if err := c.Bind(&query); err != nil {
		logger.Error(err)
		c.JSON(
			errs.Errors.ParseError.Code,
			getAllResponseError{Error: errs.Errors.ParseError.Error.Error()},
		)
		return
	}

	resp := Service.GetAll(&query)
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(
			resp.ErrResp.Code,
			getAllResponseError{Error: resp.ErrResp.Error.Error()},
		)
		return
	}

	c.JSON(200, getAllResponseSuccess{Buildings: *resp.Buildings, Count: len(*resp.Buildings)})
}

type placeTownHallBody struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Type      string `json:"type"`
	SectorID  string `json:"sectorId"`
	AccountID string `json:"accountId"`
}

// PlaceTownHall создает постройку
func (s *controller) PlaceTownHall(c *gin.Context) {
	accountID := c.GetString(`accountID`)

	// Парсинг body
	var body placeTownHallBody
	if err := tools.ParseBody(c, &body); err != nil {
		logger.Error(err)
		c.JSON(
			errs.Errors.ParseError.Code,
			createResponseError{Error: errs.Errors.ParseError.Error.Error()},
		)
		return
	}

	body.AccountID = accountID
	body.Type = `townhall`

	resp := Service.PlaceTownHall(repositories.BuildingCreate(body))
	if resp.Err != nil {
		logger.Error(resp.Err)
		c.JSON(
			resp.ErrResp.Code,
			createResponseError{Error: resp.ErrResp.Error.Error()},
		)
		return
	}

	c.JSON(200, createResponseSuccess{Building: *resp.Building})
}
