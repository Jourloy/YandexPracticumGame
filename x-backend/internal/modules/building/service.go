package building

import (
	"errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/jourloy/X-Backend/internal/cache"
	"github.com/jourloy/X-Backend/internal/config/errs"
	"github.com/jourloy/X-Backend/internal/modules/account"
	"github.com/jourloy/X-Backend/internal/modules/node"
	"github.com/jourloy/X-Backend/internal/modules/sector"
	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/storage"
	building_templates "github.com/jourloy/X-Backend/internal/templates/buildings"
)

type service struct {
	db    gorm.DB
	cache redis.Client
	nodeS node.Service
}

var Service *service

// Init создает сервис постройки
func InitService() {
	migration()

	service := service{
		db:    *storage.Database,
		cache: *cache.Client,
		nodeS: *node.InitService(),
	}

	Service = &service
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Building{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

type createResp struct {
	Err      error
	ErrResp  errs.Err
	Building *repositories.Building
}

// Create создает постройку
func (s *service) Create(create repositories.BuildingCreate) createResp {
	// Проверяем аккаунт
	accountResp := account.Service.GetOne(&repositories.AccountGet{ID: &create.AccountID})

	// Если ошибка
	if accountResp.Err != nil {
		return createResp{
			Err:     accountResp.Err,
			ErrResp: accountResp.ErrResp,
		}
	}

	// Проверяем сектор
	sectorResp := sector.Service.GetOne(&repositories.SectorGet{ID: &create.SectorID})

	// Если ошибка
	if sectorResp.Err != nil {
		return createResp{
			Err:     sectorResp.Err,
			ErrResp: sectorResp.ErrResp,
		}
	}

	// Проверяем постройки
	b := repositories.Building{
		X:         create.X,
		Y:         create.Y,
		AccountID: create.AccountID,
		SectorID:  create.SectorID,
	}
	res := s.db.First(&b, b)

	// Если ошибка
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return createResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	// Если есть
	if b.ID != `` {
		return createResp{
			Err:     errs.Errors.PlaceUsed.Error,
			ErrResp: errs.Errors.PlaceUsed,
		}
	}

	// Создание постройки
	template := building_templates.BuildingTemplates[create.Type]
	building := repositories.Building{
		ID:            uuid.NewString(),
		X:             create.X,
		Y:             create.Y,
		Type:          create.Type,
		AccountID:     create.AccountID,
		SectorID:      create.SectorID,
		MaxDurability: template.MaxDurability,
		Durability:    template.Durability,
		MaxStorage:    template.MaxStorage,
		UsedStorage:   template.UsedStorage,
		Level:         template.Level,
		AttackRange:   template.AttackRange,
		CanTrade:      template.CanTrade,
	}
	res = s.db.Create(&building)

	// Если ошибка
	if res.Error != nil {
		return createResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	return createResp{
		Building: &building,
	}
}

type getOneResp struct {
	Err      error
	ErrResp  errs.Err
	Building *repositories.Building
}

// GetOne получает постройку, подходящую под условие
func (s *service) GetOne(query *repositories.BuildingGet) getOneResp {
	building := repositories.Building{}
	res := s.db.First(&building, query)

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
		Building: &building,
	}
}

type getAllResp struct {
	Err       error
	ErrResp   errs.Err
	Buildings *[]repositories.Building
}

// GetAll получает все постройки, подходящие под условие
func (s *service) GetAll(query *repositories.BuildingGet) getAllResp {
	buildings := []repositories.Building{}
	res := s.db.First(&buildings, query)

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
		Buildings: &buildings,
	}
}

type placeTownHallResponse struct {
	Err      error
	ErrResp  errs.Err
	Building *repositories.Building
}

func (s *service) PlaceTownHall(create repositories.BuildingCreate) placeTownHallResponse {
	// Проверяем аккаунт
	accountResp := account.Service.GetOne(&repositories.AccountGet{ID: &create.AccountID})

	// Если ошибка
	if accountResp.Err != nil {
		return placeTownHallResponse{
			Err:     accountResp.Err,
			ErrResp: accountResp.ErrResp,
		}
	}

	// Проверяем сектор
	sectorResp := sector.Service.GetOne(&repositories.SectorGet{ID: &create.SectorID})

	// Если ошибка
	if sectorResp.Err != nil {
		return placeTownHallResponse{
			Err:     sectorResp.Err,
			ErrResp: sectorResp.ErrResp,
		}
	}

	// Проверяем узел
	nodeRep := s.nodeS.GetOne(&repositories.NodeGet{X: &create.X, Y: &create.Y, SectorID: &create.SectorID})

	// Если ошибка
	if sectorResp.Err != nil {
		return placeTownHallResponse{
			Err:     sectorResp.Err,
			ErrResp: sectorResp.ErrResp,
		}
	}

	// Если нельзя строить
	if !nodeRep.Node.Walkable {
		return placeTownHallResponse{
			Err:     errs.Errors.PlaceUsed.Error,
			ErrResp: errs.Errors.PlaceUsed,
		}
	}

	// Проверка наличия townhall
	townhall := repositories.Building{Type: `townhall`}
	res := s.db.First(&townhall, townhall)

	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return placeTownHallResponse{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	if townhall.ID != `` {
		return placeTownHallResponse{
			Err:     errs.Errors.TownhallExist.Error,
			ErrResp: errs.Errors.TownhallExist,
		}
	}

	// Создание townhall
	return placeTownHallResponse(s.Create(create))
}
