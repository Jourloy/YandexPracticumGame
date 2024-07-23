package item

import (
	"github.com/google/uuid"
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

// InitService создает сервис предмета
func InitService() *Service {
	migration()

	return &Service{
		db:    *storage.Database,
		cache: *cache.Client,
	}
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Item{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

type createResp struct {
	Err     error
	ErrResp errs.Err
	Item    *repositories.Item
}

// Create создает предмет
func (s *Service) Create(create repositories.ItemCreate) createResp {

	// Создаем предмет
	item := repositories.Item{
		ID:         uuid.NewString(),
		Type:       create.Type,
		X:          create.X,
		Y:          create.Y,
		ParentID:   create.ParentID,
		ParentType: create.ParentType,
		CreatorID:  create.CreatorID,
		SectorID:   create.SectorID,
	}

	s.db.Create(&item)

	return createResp{
		Item: &item,
	}
}
