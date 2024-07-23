package sector

import (
	"errors"
	"math/rand"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/jourloy/X-Backend/internal/cache"
	"github.com/jourloy/X-Backend/internal/config/errs"
	"github.com/jourloy/X-Backend/internal/modules/deposit"
	"github.com/jourloy/X-Backend/internal/modules/node"
	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/storage"
)

type service struct {
	db       gorm.DB
	cache    redis.Client
	depositS deposit.Service
	nodeS    node.Service
}

var Service *service

// InitService создает сервис сектора
func InitService() *service {
	migration()

	service := service{
		db:       *storage.Database,
		cache:    *cache.Client,
		depositS: *deposit.InitService(),
		nodeS:    *node.InitService(),
	}

	return &service
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Sector{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

type CreateOptions struct {
	// Глобальные координаты
	X int `json:"x"`
	Y int `json:"y"`

	// Насколько сложная местность. Минимум 0, максимум 100
	// TODO
	Difficult int `json:"difficult"`

	// Насколько непроходимая местность. Минимум 0, максимум 100
	// TODO
	Walkable int `json:"walkable"`

	// Обилие ресурсов. Минимум 0, максимум 100
	// TODO
	Abundance int `json:"abundance"`

	// Могут ли появится редкие ресурсы
	// TODO
	IsRare bool `json:"isRare"`
}

type createResp struct {
	Err     error
	ErrResp errs.Err
	Sector  *repositories.Sector
}

// Генерация сектора
func (s *service) Create(body CreateOptions) createResp {
	sector := repositories.Sector{
		ID: uuid.NewString(),
		X:  body.X,
		Y:  body.Y,
	}
	res := s.db.Create(&sector)
	// Если ошибка
	if res.Error != nil {
		return createResp{
			Err:     errs.Errors.DatabaseError.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	go s.generateNodes(sector.ID)

	return createResp{
		Sector: &sector,
	}
}

func (s *service) generateNodes(sectorID string) {
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			node := repositories.NodeCreate{
				X:         x,
				Y:         y,
				Walkable:  true,
				Difficult: 0,
				SectorID:  sectorID,
			}

			go s.generateDeposits(sectorID, x, y)

			if res := s.nodeS.Create(node); res.Err != nil {
				logger.Error(res.Err)
			}
		}
	}
}

func (s *service) generateDeposits(sectorID string, x int, y int) {
	resourceCreateRand := rand.Intn(10)
	if resourceCreateRand > 5 {
		resourceTypeRand := rand.Intn(2)

		resourceType := `wood`
		if resourceTypeRand == 1 {
			resourceType = `stone`
		}

		deposit := repositories.DepositCreate{
			X:        x,
			Y:        y,
			Type:     resourceType,
			SectorID: sectorID,
		}

		if res := s.depositS.Create(deposit); res.Err != nil {
			logger.Error(res.Err)
		}
	}
}

type getOneResp struct {
	Err     error
	ErrResp errs.Err
	Sector  *repositories.Sector
}

// GetOne получает сектор по id
func (s *service) GetOne(query *repositories.SectorGet) getOneResp {
	sector := repositories.Sector{}
	res := s.db.First(&sector, query)

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
		Sector: &sector,
	}
}

type getAllResp struct {
	Err     error
	ErrResp errs.Err
	Sectors *[]repositories.Sector
}

// GetAll возвращает все сектора
func (s *service) GetAll(query repositories.SectorGet) getAllResp {
	sectors := []repositories.Sector{}
	res := s.db.Find(&sectors, query)

	// Если ничего не нашли
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return getAllResp{
			Err:     errs.Errors.NotFound.Error,
			ErrResp: errs.Errors.NotFound,
		}
	}

	// Если ошибка
	if res.Error != nil {
		return getAllResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	return getAllResp{
		Sectors: &sectors,
	}
}
