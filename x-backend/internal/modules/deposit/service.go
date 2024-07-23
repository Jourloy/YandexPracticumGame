package deposit

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

// InitService создает сервис залежи
func InitService() *Service {
	migration()

	return &Service{
		db:    *storage.Database,
		cache: *cache.Client,
	}
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Deposit{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

type createResp struct {
	Err        error
	ErrMessage string
	Code       int
	Deposit    *repositories.Deposit
}

// Create создает сервис залежи
func (s *Service) Create(body repositories.DepositCreate) createResp {
	// Проверка существования сектора
	sector := repositories.Sector{}
	res := s.db.First(&sector, repositories.SectorGet{ID: &body.SectorID})

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

	// Создаем сервис залежи
	deposit := repositories.Deposit{
		ID:       uuid.NewString(),
		X:        body.X,
		Y:        body.Y,
		Type:     body.Type,
		Amount:   body.Amount,
		SectorID: body.SectorID,
	}

	res = s.db.Create(&deposit)
	if res.Error != nil {
		return createResp{
			Err:        res.Error,
			ErrMessage: `database error [deposit create]`,
			Code:       400,
		}
	}

	return createResp{
		Code:    200,
		Deposit: &deposit,
	}
}
