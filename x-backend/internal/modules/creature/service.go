package creature

import (
	"errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/jourloy/X-Backend/internal/cache"
	"github.com/jourloy/X-Backend/internal/config/errs"
	"github.com/jourloy/X-Backend/internal/modules/account"
	"github.com/jourloy/X-Backend/internal/modules/building"
	"github.com/jourloy/X-Backend/internal/modules/sector"
	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/storage"
	creature_templates "github.com/jourloy/X-Backend/internal/templates/creatures"
)

type Service struct {
	db    gorm.DB
	cache redis.Client
}

// InitService создает сервис существа
func InitService() *Service {
	migration()

	return &Service{
		db:    *storage.Database,
		cache: *cache.Client,
	}
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Creature{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

type createResp struct {
	Err      error
	ErrResp  errs.Err
	Creature *repositories.Creature
}

// Create создает существо
func (s *Service) Create(body repositories.CreatureCreate) createResp {

	// Создаем существо
	creature := repositories.Creature{
		ID:        uuid.NewString(),
		X:         body.X,
		Y:         body.Y,
		Race:      body.Race,
		IsWorker:  body.IsWorker,
		IsWarrior: body.IsWarrior,
		IsTrader:  body.IsTrader,
		AccountID: body.AccountID,
		SectorID:  body.SectorID,
	}

	// Шаблон
	template := creature_templates.CreatureTemplate[body.Race]

	requireFood := template.RequireFood
	if body.IsWorker {
		requireFood += 0.2
	}
	if body.IsTrader {
		requireFood += 0.2
	}
	if body.IsWarrior {
		requireFood += 0.4
	}

	creature.MaxStorage = template.MaxStorage
	creature.UsedStorage = template.UsedStorage
	creature.FatiguePerStep = template.FatiguePerStep
	creature.Fatigue = template.Fatigue
	creature.MaxHealth = template.MaxHealth
	creature.Health = template.Health
	creature.RequireFood = requireFood

	res := s.db.Create(&creature)
	if res.Error != nil {
		return createResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	return createResp{
		Creature: &creature,
	}
}

type spawnCreature struct {
	IsWorker  bool   `form:"isWorker" json:"isWorker"`
	IsTrader  bool   `form:"IsTrader" json:"IsTrader"`
	IsWarrior bool   `form:"IsWarrior" json:"IsWarrior"`
	SectorID  string `form:"sectorId" json:"sectorId"`
	AccountID string
}

type spawnResponse struct {
	Err      error
	ErrResp  errs.Err
	Creature *repositories.Creature
}

// Создает существо в townhall
func (s *Service) Spawn(spawn spawnCreature) spawnResponse {
	// Проверяем аккаунт
	accountResp := account.Service.GetOne(&repositories.AccountGet{ID: &spawn.AccountID})

	// Если ошибка
	if accountResp.Err != nil {
		return spawnResponse{
			Err:     accountResp.Err,
			ErrResp: accountResp.ErrResp,
		}
	}

	// Проверяем сектор
	sectorRes := sector.Service.GetOne(&repositories.SectorGet{ID: &spawn.SectorID})

	// Если ошибка
	if sectorRes.Err != nil {
		return spawnResponse{
			Err:     sectorRes.Err,
			ErrResp: sectorRes.ErrResp,
		}
	}

	// Проверяем townhall
	bType := `townhall`
	bSectorID := sectorRes.Sector.ID
	bAccountID := spawn.AccountID
	buildingRes := building.Service.GetOne(&repositories.BuildingGet{Type: &bType, SectorID: &bSectorID, AccountID: &bAccountID})

	// Если ошибка
	if buildingRes.Err != nil {
		return spawnResponse{
			Err:     buildingRes.Err,
			ErrResp: buildingRes.ErrResp,
		}
	}

	// Подсчет стоимости существа
	price := 0
	if spawn.IsTrader {
		price += 60
	}
	if spawn.IsWarrior {
		price += 80
	}
	if spawn.IsTrader {
		price += 40
	}

	// Если не хватает денег
	if accountResp.Account.Balance-price < 0 {
		return spawnResponse{
			Err:     errs.NotEnoughBalance.Error,
			ErrResp: errs.NotEnoughBalance,
		}
	}

	// Создание существа
	creature := repositories.CreatureCreate{
		X:         buildingRes.Building.X,
		Y:         buildingRes.Building.Y,
		Race:      accountResp.Account.Race,
		IsWorker:  spawn.IsWorker,
		IsTrader:  spawn.IsTrader,
		IsWarrior: spawn.IsWarrior,
		SectorID:  spawn.SectorID,
		AccountID: spawn.AccountID,
	}

	res := s.Create(creature)
	if res.Err != nil {
		return spawnResponse{
			Err:     res.Err,
			ErrResp: res.ErrResp,
		}
	}

	// Обновление баланса
	accountResp.Account.Balance -= price
	account.Service.UpdateOne(*accountResp.Account)

	return spawnResponse{
		Creature: res.Creature,
	}
}

type getOneResp struct {
	Err      error
	ErrResp  errs.Err
	Creature *repositories.Creature
}

// GetOne получает первое сущетсво, подходящее под условие
func (s *Service) GetOne(query *repositories.CreatureGet) getOneResp {
	creature := repositories.Creature{}
	res := s.db.First(&creature, query)

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
		Creature: &creature,
	}
}

type getAllResp struct {
	Err       error
	ErrResp   errs.Err
	Creatures *[]repositories.Creature
}

// GetAll получает всех существ, подходящих под условие
func (s *Service) GetAll(query *repositories.CreatureGet) getAllResp {
	creatures := []repositories.Creature{}
	res := s.db.Find(&creatures, query)

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
		Creatures: &creatures,
	}
}

type moveCreature struct {
	ByX        int    `form:"byX" json:"byX"`
	ByY        int    `form:"byY" json:"byY"`
	CreatureID string `form:"creatureId" json:"creatureId"`
	AccountID  string
}

type moveCreatureRes struct {
	Err      error
	ErrResp  errs.Err
	Creature *repositories.Creature
}

// Передвигает существо
func (s *Service) Move(move moveCreature) moveCreatureRes {
	// Проверяем аккаунт
	accountRes := account.Service.GetOne(&repositories.AccountGet{ID: &move.AccountID})

	// Если ошибка
	if accountRes.Err != nil {
		return moveCreatureRes{
			Err:     accountRes.Err,
			ErrResp: accountRes.ErrResp,
		}
	}

	// Проверяем существо
	creatureRes := s.GetOne(&repositories.CreatureGet{ID: &move.CreatureID})

	// Если ошибка
	if creatureRes.Err != nil {
		return moveCreatureRes{
			Err:     creatureRes.Err,
			ErrResp: creatureRes.ErrResp,
		}
	}

	// Можно ходить только на одну клетку
	// TODO усталость
	if move.ByX > 0 {
		move.ByX = 1
	}
	if move.ByX < 0 {
		move.ByX = -1
	}
	if move.ByY > 0 {
		move.ByY = 1
	}
	if move.ByY < 0 {
		move.ByY = -1
	}

	creature := *creatureRes.Creature
	creature.X += move.ByX
	creature.Y += move.ByY

	res := s.UpdateOne(creature)
	// Если ошибка
	if creatureRes.Err != nil {
		return moveCreatureRes{
			Err:     res.Err,
			ErrResp: res.ErrResp,
		}
	}

	return moveCreatureRes{
		Creature: &creature,
	}
}

type updateOneResp struct {
	Err      error
	ErrResp  errs.Err
	Creature *repositories.Creature
}

func (s *Service) UpdateOne(creature repositories.Creature) updateOneResp {
	res := s.db.Save(creature)

	// Если ошибка
	if res.Error != nil {
		return updateOneResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	return updateOneResp{
		Creature: &creature,
	}
}
