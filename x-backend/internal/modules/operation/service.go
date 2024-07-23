package operation

import (
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/jourloy/X-Backend/internal/cache"
	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/storage"
)

type Service struct {
	db    gorm.DB
	cache redis.Client
}

// Init создает сервис операций
func InitService() *Service {
	migration()

	return &Service{
		db:    *storage.Database,
		cache: *cache.Client,
	}
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Operation{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

type createResp struct {
	Err        error
	ErrMessage string
	Code       int
	Operation  *repositories.Operation
}

// Create создает операцию
func (s *Service) Create(body repositories.OperationCreate) createResp {
	// Проверка аккаунта
	account := repositories.Account{}
	res := s.db.First(&account, repositories.AccountGet{ID: &body.AccountID})

	// Если ничего не нашли
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return createResp{
			Err:        errors.New(`account not found`),
			ErrMessage: `account not found`,
			Code:       404,
		}
	}

	// Если ошибка
	if res.Error != nil {
		return createResp{
			Err:        res.Error,
			ErrMessage: `database error [account find]`,
			Code:       404,
		}
	}

	// Проверка сектора
	sector := repositories.Sector{}
	res = s.db.First(&sector, repositories.SectorGet{ID: &body.SectorID})

	// Если ничего не нашли
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return createResp{
			Err:        errors.New(`sector not found`),
			ErrMessage: `sector not found`,
			Code:       404,
		}
	}

	// Если ошибка
	if res.Error != nil {
		return createResp{
			Err:        res.Error,
			ErrMessage: `database error [sector find]`,
			Code:       404,
		}
	}

	// Проверка постройки
	b := repositories.Building{
		ID: body.BuildingID,
	}
	res = s.db.First(&b, b)

	// Если ничего не нашли
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return createResp{
			Err:        errors.New(`building not found`),
			ErrMessage: `building not found`,
			Code:       404,
		}
	}

	// Если ошибка
	if res.Error != nil {
		return createResp{
			Err:        res.Error,
			ErrMessage: `database error [building create]`,
			Code:       404,
		}
	}

	// Создание операции
	operation := repositories.Operation{
		Price:      body.Price,
		Type:       *body.Type,
		Amount:     body.Amount,
		Name:       body.Name,
		IsResource: body.IsResource,
		IsItem:     body.IsItem,
		BuildingID: body.BuildingID,
		SectorID:   body.SectorID,
		AccountID:  body.AccountID,
	}

	res = s.db.Create(&operation)
	if res.Error != nil {
		return createResp{
			Err:        res.Error,
			ErrMessage: `database error [operation create]`,
			Code:       400,
		}
	}

	return createResp{
		Code:      200,
		Operation: &operation,
	}
}
