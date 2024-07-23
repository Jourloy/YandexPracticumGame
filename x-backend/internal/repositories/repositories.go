package repositories

import (
	"time"

	"gorm.io/gorm"
)

// Модель аккаунта
type Account struct {
	// Задается при создании

	Race     string `json:"race"`
	Username string `json:"username"`

	// Задается по умолчанию

	ApiKey  string `json:"apiKey"`
	Balance int    `json:"balance"`
	IsAdmin bool   `json:"-"`

	// Мета

	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type AccountCreate struct {
	Username string `form:"username" json:"username"`
	Race     string `form:"race" json:"race"`
}

type AccountGet struct {
	ID       *string `form:"id" json:"id"`
	ApiKey   *string `form:"apiKey" json:"apiKey"`
	Username *string `form:"username" json:"username"`
	Balance  *int    `form:"balance" json:"balance"`
	Race     *string `form:"race" json:"race"`
	Limit    *int    `form:"limit" json:"limit"`
}

type AccountRepository interface {
	Create(create *AccountCreate) (*Account, error)
	GetOne(query *AccountGet) (*Account, error)
	UpdateOne(account *Account) error
	DeleteOne(account *Account) error
}

// Модель сектора
type Sector struct {
	// Задается при создании

	X int `json:"x"`
	Y int `json:"y"`

	// Дети

	Nodes     []Node     `json:"nodes"`
	Buildings []Building `json:"buildings"`
	Plans     []Plan     `json:"plans"`
	Creatures []Creature `json:"creatures"`
	Deposits  []Deposit  `json:"deposits"`
	Resources []Resource `json:"resources" gorm:"foreignKey:ParentID"`
	Items     []Item     `json:"items" gorm:"foreignKey:ParentID"`

	// Мета
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type SectorCreate struct {
	X int `form:"x" json:"x"`
	Y int `form:"y" json:"y"`
}

// Структура поиска сектора
type SectorGet struct {
	ID    *string `form:"id" json:"id"`
	X     *int    `form:"x" json:"x"`
	Y     *int    `form:"y" json:"y"`
	Limit *int    `form:"limit" json:"limit"`
}

// Репозиторий сектора
type SectorRepository interface {
	Create(sector *SectorCreate) (*Sector, error)
	GetOne(query *SectorGet) (*Sector, error)
	GetAll(query *SectorGet) (*[]Sector, error)
	UpdateOne(sector *Sector) error
	DeleteOne(sector *Sector) error
}

// Модель узла
type Node struct {
	// Задается при создании

	X         int  `json:"x"`
	Y         int  `json:"y"`
	Walkable  bool `json:"walkable"`
	Difficult int  `json:"difficult"`

	// Родители

	SectorID string `json:"sectorId"`

	// Мета
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type NodeCreate struct {
	X         int    `form:"x" json:"x"`
	Y         int    `form:"y" json:"y"`
	Walkable  bool   `form:"walkable" json:"walkable"`
	Difficult int    `form:"difficult" json:"difficult"`
	SectorID  string `form:"sectorId" json:"sectorId"`
}

// Модель поиска узла
type NodeGet struct {
	X         *int    `form:"x" json:"x"`
	Y         *int    `form:"y" json:"y"`
	Walkable  *bool   `form:"walkable" json:"walkable"`
	Difficult *int    `form:"difficult" json:"difficult"`
	SectorID  *string `form:"sectorId" json:"sectorId"`
	Limit     *int    `form:"limit" json:"limit"`
}

// Репозиторий сектора
type NodeRepository interface {
	Create(create *NodeCreate) (*Node, error)
	GetOne(query *NodeGet) (*Node, error)
	GetAll(query *NodeGet) (*[]Node, error)
	UpdateOne(node *Node) error
	DeleteOne(node *Node) error
}

// Модель залежи
type Deposit struct {
	// Задается при создании

	X      int    `json:"x"`
	Y      int    `json:"y"`
	Type   string `json:"type"`
	Amount int    `json:"amount"`

	// Родители

	SectorID string `json:"sectorId"`

	// Мета

	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Структура создания залежи
type DepositCreate struct {
	X        int    `form:"x" json:"x"`
	Y        int    `form:"y" json:"y"`
	Type     string `form:"type" json:"type"`
	Amount   int    `form:"amount" json:"amount"`
	SectorID string `form:"sectorId" json:"sectorId"`
}

// Структура поиска залежей
type DepositGet struct {
	Type     *string `form:"type" json:"type"`
	Amount   *int    `form:"amount" json:"amount"`
	X        *int    `form:"x" json:"x"`
	Y        *int    `form:"y" json:"y"`
	SectorID *string `form:"sectorId" json:"sectorId"`
	Limit    *int    `form:"limit" json:"limit"`
}

// Репозиторий залежей
type DepositRepository interface {
	Create(create DepositCreate) (*Deposit, error)
	GetOne(query DepositGet) (*Deposit, error)
	GetAll(query DepositGet) (*[]Deposit, error)
	UpdateOne(deposit *Deposit) error
	DeleteOne(deposit *Deposit) error
}

// Модель ресурсов
type Resource struct {
	// Задается при создании

	X    int    `json:"x"`
	Y    int    `json:"y"`
	Type string `json:"type"`

	// Задается по умолчанию

	Amount int `json:"amount"`

	// Задается по шаблону

	Weight int `json:"weight"`

	// Родители

	ParentID   string `json:"parentId"`
	ParentType string `json:"parentType"`
	SectorID   string `json:"sectorId"`
	CreatorID  string `json:"creatorId"`

	// Мета

	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Структура создания ресурсов
type ResourceCreate struct {
	X          int    `form:"x" json:"x"`
	Y          int    `form:"y" json:"y"`
	Type       string `form:"type" json:"type"`
	ParentID   string `json:"parentId"`
	ParentType string `json:"parentType"`
	SectorID   string `json:"sectorId"`
	CreatorID  string `json:"creatorId"`
}

// Структура поиска ресурсов
type ResourceGet struct {
	Type       *string `form:"type" json:"type"`
	Amount     *int    `form:"amount" json:"amount"`
	Weight     *int    `form:"weight" json:"weight"`
	X          *int    `form:"x" json:"x"`
	Y          *int    `form:"y" json:"y"`
	ParentID   *string `form:"parentId" json:"parentId"`
	ParentType *string `form:"parentType" json:"parentType"`
	SectorID   *string `form:"sectorId" json:"sectorId"`
	CreatorID  *string `form:"creatorId" json:"creatorId"`
	Limit      *int    `form:"limit" json:"limit"`
}

// Репозиторий ресурсов
type ResourceRepository interface {
	Create(create ResourceCreate) (*Resource, error)
	GetOne(query ResourceGet) (*Resource, error)
	GetAll(query ResourceGet) (*[]Resource, error)
	UpdateOne(resource *Resource)
	DeleteOne(resource *Resource)
}

// Модель предмета
type Item struct {
	// Задается при создании

	X    int    `json:"x"`
	Y    int    `json:"y"`
	Type string `json:"type"`

	// Задается по шаблону

	Weight int `json:"weight"`

	// Родители

	ParentID   string `json:"parentId"`
	ParentType string `json:"parentType"`
	SectorID   string `json:"sectorId"`
	CreatorID  string `json:"creatorId"`

	// Мета

	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Структура создания предмета
type ItemCreate struct {
	X          int    `form:"Y" json:"X"`
	Y          int    `form:"Y" json:"Y"`
	Type       string `form:"Type" json:"string"`
	ParentID   string `form:"parentId" json:"parentId"`
	ParentType string `form:"parentType" json:"parentType"`
	SectorID   string `form:"sectorId" json:"sectorId"`
	CreatorID  string `form:"creatorId" json:"creatorId"`
}

// Структура поиска предмета
type ItemGet struct {
	ID         *string `form:"id" json:"id"`
	Type       *string `form:"type" json:"type"`
	X          *int    `form:"x" json:"x"`
	Y          *int    `form:"y" json:"y"`
	ParentID   *string `form:"parentId" json:"parentId"`
	ParentType *string `form:"parentType" json:"parentType"`
	SectorID   *string `form:"sectorID" json:"sectorID"`
	CreatorID  *string `form:"creatorID" json:"creatorID"`
	Limit      *int    `form:"limit" json:"limit"`
}

// Репозиторий предмета
type ItemRepository interface {
	Create(item *Item)
	GetOne(item *Item)
	GetAll(query ItemGet) []Item
	UpdateOne(item *Item)
	DeleteOne(item *Item)
}

// Модель операции
type Operation struct {
	// Динамические поля, задаются пользователем
	Price      int    `json:"price"`
	Amount     int    `json:"amount"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	IsResource bool   `json:"isResource"`
	IsItem     bool   `json:"isItem"`

	// Родители
	BuildingID string `json:"buildingID"`
	SectorID   string `json:"sectorId"`
	AccountID  string `json:"accountId"`

	// Мета
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Структура создания операции
type OperationCreate struct {
	Price      int     `form:"price" json:"price"`
	Amount     int     `form:"amount" json:"amount"`
	Type       *string `form:"type" json:"type"`
	Name       string  `form:"name" json:"name"`
	IsResource bool    `form:"isResource" json:"isResource"`
	IsItem     bool    `form:"isItem" json:"isItem"`
	BuildingID string  `form:"buildingID" json:"buildingID"`
	SectorID   string  `form:"sectorID" json:"sectorID"`
	AccountID  string  `form:"accountID" json:"accountID"`
}

// Структура поиска операции
type OperationGet struct {
	ID         *string  `form:"id" json:"id,omitempty"`
	Price      *int     `form:"price" json:"price,omitempty"`
	Amount     *int     `form:"amount" json:"amount,omitempty"`
	Type       *string  `form:"type" json:"type,omitempty"`
	Name       *string  `form:"name" json:"name,omitempty"`
	IsResource *bool    `form:"isResource" json:"isResource,omitempty"`
	IsItem     *bool    `form:"isItem" json:"isItem,omitempty"`
	BuildingID *float64 `form:"buildingID" json:"buildingID,omitempty"`
	SectorID   *string  `form:"sectorID" json:"sectorId,omitempty"`
	AccountID  *string  `form:"accountID" json:"accountId,omitempty"`
	Limit      *int     `form:"limit" json:"limit,omitempty"`
}

// Репозиторий операции
type OperationRepository interface {
	Create(create *OperationCreate) (*Operation, error)
	GetOne(query *OperationGet) (*Operation, error)
	GetAll(query *OperationGet) (*[]Operation, error)
	UpdateOne(operation *Operation) error
	DeleteOne(operation *Operation) error
}

//////// Постройки ////////

// Модель постройки
type Building struct {
	ID string `json:"id" gorm:"primarykey"`

	X int `json:"x"`
	Y int `json:"y"`

	// Динамические поля, задаются пользователем
	Type string `json:"type"`

	// Динамические поля, задаются шаблоном
	MaxDurability int  `json:"maxDurability"`
	Durability    int  `json:"durability"`
	MaxStorage    int  `json:"maxStorage"`
	UsedStorage   int  `json:"usedStorage"`
	Level         int  `json:"level"`
	AttackRange   int  `json:"attackRange"`
	CanTrade      bool `json:"canTrade"`

	// Дети
	Items      []Item      `json:"items" gorm:"foreignKey:ParentID"`
	Resources  []Resource  `json:"resources" gorm:"foreignKey:ParentID"`
	Operations []Operation `json:"operations"`

	// Родители
	SectorID  string `json:"sectorId"`
	AccountID string `json:"accountId"`

	// Мета
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Структура создания постройки
type BuildingCreate struct {
	X         int    `form:"x" json:"x"`
	Y         int    `form:"y" json:"y"`
	Type      string `form:"type" json:"type"`
	SectorID  string `form:"sectorID" json:"sectorId"`
	AccountID string `form:"accountID" json:"accountId"`
}

// Структура поиска постройки
type BuildingGet struct {
	ID            *string  `form:"id" json:"id,omitempty"`
	Type          *string  `form:"type" json:"type,omitempty"`
	MaxDurability *int     `form:"maxDurability" json:"maxDurability,omitempty"`
	Durability    *int     `form:"durability" json:"durability,omitempty"`
	MaxStorage    *float64 `form:"maxStorage" json:"maxStorage,omitempty"`
	UsedStorage   *float64 `form:"usedStorage" json:"usedStorage,omitempty"`
	Level         *float64 `form:"level" json:"level,omitempty"`
	AttackRange   *float64 `form:"attackRange" json:"attackRange,omitempty"`
	CanTrade      *bool    `form:"canTrade" json:"canTrade,omitempty"`
	SectorID      *string  `form:"sectorId" json:"sectorId,omitempty"`
	AccountID     *string  `form:"accountId" json:"accountId,omitempty"`
	Limit         *int     `form:"limit" json:"limit,omitempty"`
}

// Репозиторий постройки
type BuildingRepository interface {
	Create(create *BuildingCreate) (*Building, error)
	GetOne(query *BuildingGet) (*Building, error)
	GetAll(query *BuildingGet) (*[]Building, error)
	UpdateOne(building *Building) error
	DeleteOne(building *Building) error
}

// Модель планируемой постройки
type Plan struct {
	ID string `json:"id" gorm:"primarykey"`

	X int `json:"x"`
	Y int `json:"y"`

	// Динамические поля, задаются пользователем
	Type string `json:"type"`

	MaxProgress int `json:"maxProgress"`
	Progress    int `json:"progress"`

	// Дети
	Items     []Item     `json:"items" gorm:"foreignKey:ParentID"`
	Resources []Resource `json:"resources" gorm:"foreignKey:ParentID"`

	// Родители
	SectorID  string `json:"sectorId"`
	AccountID string `json:"accountId"`

	// Мета
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Структура создания планируемой постройки
type PlanCreate struct {
	X         int    `form:"x" json:"x"`
	Y         int    `form:"y" json:"y"`
	Type      string `form:"type" json:"type"`
	SectorID  string `form:"sectorId" json:"sectorId"`
	AccountID string `form:"accountId" json:"accountId"`
}

// Структура поиска планируемой постройки
type PlanGetAll struct {
	MaxProgress *int    `form:"maxProgress" json:"maxProgress"`
	Progress    *int    `form:"progress" json:"progress"`
	Type        *string `form:"type" json:"type"`
	Y           *int    `form:"y" json:"y"`
	X           *int    `form:"x" json:"x"`
	Limit       *int    `form:"limit" json:"limit"`
}

// Репозиторий планируемой постройки
type IPlanRepository interface {
	Create(plan *PlanCreate)
	GetOne(plan *Plan)
	GetAll(query PlanGetAll, accountID string) []Plan
	UpdateOne(plan *Plan)
	DeleteOne(plan *Plan)
}

//////// Существа ////////

// Модель существа
type Creature struct {
	ID string `json:"id" gorm:"primarykey"`

	X int `json:"x"`
	Y int `json:"y"`

	// Динамические поля, задаются пользователем
	Race      string `json:"race"`
	IsWorker  bool   `json:"isWorker"`
	IsTrader  bool   `json:"isTrader"`
	IsWarrior bool   `json:"isWarrior"`

	// Динамические поля, задаются шаблоном
	MaxStorage         int     `json:"maxStorage"`
	UsedStorage        int     `json:"usedStorage"`
	RequireFood        float64 `json:"requireFood"`
	FatiguePerStep     float64 `json:"fatiguePerStep"`
	FatigueModificator float64 `json:"fatigueModificator"`
	Fatigue            float64 `json:"fatigue"`
	MaxHealth          int     `json:"maxHealth"`
	Health             int     `json:"health"`

	// Дети
	Items     []Item     `json:"items" gorm:"foreignKey:ParentID"`
	Resources []Resource `json:"resources" gorm:"foreignKey:ParentID"`

	// Родители
	SectorID  string `json:"sectorId"`
	AccountID string `json:"accountId"`

	// Мета
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Структура создания существа
type CreatureCreate struct {
	X         int    `form:"x" json:"x"`
	Y         int    `form:"y" json:"y"`
	Race      string `form:"race" json:"race"`
	IsWorker  bool   `form:"isWorker" json:"isWorker"`
	IsTrader  bool   `form:"isTrader" json:"isTrader"`
	IsWarrior bool   `form:"isWarrior" json:"isWarrior"`
	SectorID  string `form:"sectorId" json:"sectorId"`
	AccountID string `form:"accountId" json:"accountId"`
}

// Структура поиска существа
type CreatureGet struct {
	ID                 *string  `form:"id" json:"id,omitempty"`
	Race               *string  `form:"race" json:"race,omitempty"`
	MaxStorage         *int     `form:"maxStorage" json:"maxStorage,omitempty"`
	UsedStorage        *int     `form:"usedStorage" json:"usedStorage,omitempty"`
	RequireCoins       *float64 `form:"requireCoins" json:"requireCoins,omitempty"`
	RequireFood        *float64 `form:"requireFood" json:"requireFood,omitempty"`
	Fatigue            *float64 `form:"fatigue" json:"fatigue,omitempty"`
	FatiguePerStep     *float64 `form:"fatiguePerStep" json:"fatiguePerStep,omitempty"`
	FatigueModificator *float64 `form:"fatigueModificator" json:"fatigueModificator,omitempty"`
	MaxHealth          *int     `form:"maxHealth" json:"maxHealth,omitempty"`
	Health             *int     `form:"health" json:"health,omitempty"`
	IsWorker           *bool    `form:"isWorker" json:"isWorker,omitempty"`
	IsTrader           *bool    `form:"isTrader" json:"isTrader,omitempty"`
	IsWarrior          *bool    `form:"isWarrior" json:"isWarrior,omitempty"`
	SectorID           *string  `form:"sectorId" json:"sectorId,omitempty"`
	AccountID          *string  `form:"accountId" json:"accountId,omitempty"`
	Limit              *int     `form:"limit" json:"limit,omitempty"`
}

// Репозиторий существа
type CreatureRepository interface {
	Create(creatrue *CreatureCreate) (*Creature, error)
	GetOne(query *CreatureGet) (*Creature, error)
	GetAll(query *CreatureGet) (*[]Creature, error)
	UpdateOne(creature *Creature) error
	DeleteOne(creature *Creature) error
}
