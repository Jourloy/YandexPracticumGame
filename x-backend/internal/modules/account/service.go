package account

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/jourloy/X-Backend/internal/cache"
	"github.com/jourloy/X-Backend/internal/config/errs"
	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/storage"
	"github.com/jourloy/X-Backend/internal/tools"
)

type service struct {
	done  chan struct{}
	db    gorm.DB
	cache redis.Client
}

var Service *service

// InitService создает сервис аккаунта
func initService() {
	migration()

	service := service{
		done:  make(chan struct{}),
		db:    *storage.Database,
		cache: *cache.Client,
	}

	go service.tickers()

	Service = &service
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&repositories.Account{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

func (s *service) tickers() {
	addMoneyTicker := time.NewTicker(1 * time.Second)
	defer addMoneyTicker.Stop()

	for {
		select {
		case <-s.done:
			return
		case <-addMoneyTicker.C:
			s.addMoney()
		}
	}
}

// addMoney добавляет пользователю по 1 монете
func (s *service) addMoney() {
	// Ищем townhalls
	townhalls := []repositories.Building{}
	res := s.db.Find(&townhalls, repositories.Building{Type: `townhall`})

	// Если ничего нет
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return
	}

	// Если ошибка
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		logger.Error(`Add money ticker`, `error`, res.Error)
	}

	// ID аккаунтов пользователей
	ids := []string{}
	for _, t := range townhalls {
		if !tools.Contains(ids, t.AccountID) {
			ids = append(ids, t.AccountID)
		}
	}

	// Ищем аккаунты
	accounts := []repositories.Account{}
	res = s.db.Find(&accounts, ids)

	// Если ничего нет
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return
	}

	// Если ошибка
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		logger.Error(`Add money ticker`, `error`, res.Error)
	}

	for _, a := range accounts {
		a.Balance += 1
		go s.db.Save(a)
	}
}

type createResp struct {
	Err     error
	ErrResp errs.Err
	Account *repositories.Account
}

// Create создает аккаунт
func (s *service) Create(body repositories.AccountCreate) createResp {
	account := repositories.Account{}
	res := s.db.First(&account, repositories.Account{Username: body.Username})

	// Если ошибка
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return createResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	// Если есть
	if (account != repositories.Account{}) {
		return createResp{
			Err:     errs.Errors.AlreadyExist.Error,
			ErrResp: errs.Errors.AlreadyExist,
		}
	}

	// Создаем аккаунт
	account = repositories.Account{
		ID:       uuid.NewString(),
		ApiKey:   uuid.NewString(),
		Username: body.Username,
		Race:     body.Race,
		Balance:  0,
		IsAdmin:  false,
	}

	res = s.db.Create(&account)
	if res.Error != nil {
		return createResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	return createResp{
		Account: &account,
	}
}

type getOneResp struct {
	Err     error
	ErrResp errs.Err
	Account *repositories.Account
}

// GetOne получает первый аккаунт, попавший под условие
func (s *service) GetOne(query *repositories.AccountGet) getOneResp {
	account := repositories.Account{}
	res := s.db.First(&account, query)

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
		Account: &account,
	}
}

type updateOneResp struct {
	Err     error
	ErrResp errs.Err
	Account *repositories.Account
}

func (s *service) UpdateOne(account repositories.Account) updateOneResp {
	res := s.db.Save(account)

	// Если ошибка
	if res.Error != nil {
		return updateOneResp{
			Err:     res.Error,
			ErrResp: errs.Errors.DatabaseError,
		}
	}

	return updateOneResp{
		Account: &account,
	}
}
