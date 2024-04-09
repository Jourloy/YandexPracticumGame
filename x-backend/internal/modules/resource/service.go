package resource

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/jourloy/X-Backend/internal/cache"
	"github.com/jourloy/X-Backend/internal/config/errs"
	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/storage"
)

type Service struct {
	db    gorm.DB
	cache redis.Client
}

// InitService создает сервис ресурса
func InitService() *Service {
	migration()

	return &Service{
		db:    *storage.Database,
		cache: *cache.Client,
	}
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Resource{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

type createResp struct {
	Err      error
	ErrResp  errs.Err
	Resource repositories.Resource
}

// Create создает ресурс
func (s *Service) Create(create repositories.ResourceCreate) createResp {

	resource := repositories.Resource{
		X:          create.X,
		Y:          create.Y,
		Type:       create.Type,
		ParentID:   create.ParentID,
		ParentType: create.ParentType,
		SectorID:   create.SectorID,
		CreatorID:  create.CreatorID,
	}

	res := s.db.Create(&resource)
	// Если ошибка
	if res.Error != nil {
		return createResp{
			Err:     errs.Errors.DatabaseError.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	return createResp{Resource: resource}
}
