package plan

import (
	"errors"

	"github.com/google/uuid"
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

// InitService создает сервис планируемой постройки
func InitService() *Service {
	migration()

	return &Service{
		db:    *storage.Database,
		cache: *cache.Client,
	}
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Plan{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

type createResp struct {
	Err        error
	ErrMessage string
	Code       int
	Plan       *repositories.Plan
}

// Create создает планируемую постройку
func (s *Service) Create(create repositories.PlanCreate) createResp {
	// Проверка аккаунта
	account := repositories.Account{}
	res := s.db.First(&account, repositories.AccountGet{ID: &create.AccountID})

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
	res = s.db.First(&sector, repositories.SectorGet{ID: &create.SectorID})

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

	plan := repositories.Plan{
		ID:          uuid.NewString(),
		MaxProgress: 100,
		Progress:    0,
		X:           create.X,
		Y:           create.Y,
		Type:        create.Type,
		AccountID:   create.AccountID,
	}

	res = s.db.Create(&plan)

	// Если ошибка
	if res.Error != nil {
		return createResp{
			Err:        res.Error,
			ErrMessage: `database error [building create]`,
			Code:       404,
		}
	}

	return createResp{
		Plan: &plan,
	}
}
