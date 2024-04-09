package node

import (
	"errors"

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

// InitService создает сервис аккаунта
func InitService() *Service {
	migration()

	return &Service{
		db:    *storage.Database,
		cache: *cache.Client,
	}
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Node{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

type createResp struct {
	Err     error
	ErrResp errs.Err
	Node    *repositories.Node
}

// Create создает узел
func (s *Service) Create(create repositories.NodeCreate) createResp {
	node := repositories.Node{}
	res := s.db.First(
		&node,
		repositories.Node{
			SectorID: create.SectorID,
			X:        create.X,
			Y:        create.Y,
		},
	)

	// Если ошибка
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return createResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	// Если есть
	if node.ID != `` {
		return createResp{
			Err:     errs.Errors.AlreadyExist.Error,
			ErrResp: errs.Errors.AlreadyExist,
		}
	}

	// Создаем узел
	node = repositories.Node{
		ID:        uuid.NewString(),
		SectorID:  create.SectorID,
		X:         create.X,
		Y:         create.Y,
		Walkable:  create.Walkable,
		Difficult: create.Difficult,
	}

	res = s.db.Create(&node)
	if res.Error != nil {
		return createResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	return createResp{
		Node: &node,
	}
}

type getOneResp struct {
	Err     error
	ErrResp errs.Err
	Node    *repositories.Node
}

// GetOne получает первый узел, попавший под условие
func (s *Service) GetOne(query *repositories.NodeGet) getOneResp {
	node := repositories.Node{}
	res := s.db.First(&node, query)

	// Если ничего не нашли
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return getOneResp{
			Err:     errs.Errors.NotFound.Error,
			ErrResp: errs.Errors.NotFound,
		}
	}

	// Если ошибка
	if res.Error != nil {
		return getOneResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	return getOneResp{
		Node: &node,
	}
}
