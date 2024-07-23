# Как добавить модель

1. Инициализация модели в `repositories.go`

Всегда любым моделям нужно добавлять 4 поля: `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`

```golang
// Модель примера
type Example struct {
	// Задается при создании

	Name string `json:"name"`

	// Задается по умолчанию

	DefaultValue string `json:"-"`

	// Дети

	Childs []Child `json:"childs"`

	// Родители

	ParentID string `json:"parentId"`

	// Мета

	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

Модель обязательно должна иметь структуру для создания. В нее помещаются все поля, которые задаются при создании и родители модели

```golang
type ExampleCreate struct {
	Name      string `json:"name"`
	ExampleID string `json:"exmapleId"`
	ParentID  string `json:"parentId"`
}
```

У модели обязательно должна быть структура для поиска. В нее добавляется `Limit`, а также все поля, где `json` не равен `-`. Все типы должны быть с `*`

```golang
// Структура поиска примера
type ExampleGet struct {
	ID *string
	Name *string
	ExampleID *string
	ParentID *string
	Limit *int
}
```

У модели обязательно должна быть структура репозитория. Все функции должны возвращать ошибку

```golang
// Репозиторий примера
type ExampleRepository interface {
	// Create создает модель
	Create(example *ExampleCreate) (*Example, error)
	// GetOne возвращает первую модель, попавший под условие
	GetOne(query *ExampleGet) (*Example, error)
	// GetAll возвращает все модели, попавшие под условие
	GetAll(query *ExampleGet) (*[]Example, error)
	// UpdateOne обновляет модель
	UpdateOne(example *Example) error
	// DeleteOne удаляет модель
	DeleteOne(example *Example) error
}
```

2. Добавление репозитория

В папке `repositories` необходимо создать папку `имя_модели_rep` (`example_rep`),
а в ней аналогичный файл с расширением .go

Файл должен содержать логгер

```golang
var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[database-example]`,
		Level:  log.DebugLevel,
	})
)
```

Определение репозитория

```golang
var Repository repositories.ExampleRepository
```

Структуру репозитория

```golang
type ExampleRepository struct {
	db gorm.DB
}
```

И функцию инициализации вместе с автомиграцией

```golang
// Init создает репозиторий примера
func Init() {
	go migration()

	Repository = &ExampleRepository{
		db: *storage.Database,
	}
}

func migration() {
	if err := Database.AutoMigrate(
		&repositories.Example{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}
```

## Основные функции репозитория

`ExampleRepository` **должен** содержать все функции, как и его интрефейс `repositories.ExampleRepository`. Также они должны принимать и возвращать такие же параметры

### Create

```golang
func (r *ExampleRepository) Create(create *repositories.ExampleCreate) (*repositories.Example, error) {
	// Проверка, есть ли уже такой example
	e, err := r.GetOne(&repositories.ExampleGet{Name: &create.Name})
	if err != nil {
		return nil, err
	}

	// Если есть
	if e != nil {
		return nil, errors.New(`example already exist`)
	}

	// Создаем example
	example := repositories.Example{
		ID:           uuid.NewString(),
		Name:         create.Name,
		DefaultValue: `default`,
	}

	res := r.db.Create(&example)
	if res.Error != nil {
		return nil, res.Error
	}

	return &example, nil
}
```

В начала идет проверка, если модель должна быть уникальной (например, аккаунт. В системе может быть только один аккаунт с определенным именем)

```golang
// Проверка, есть ли уже такой example
e, err := r.GetOne(&repositories.ExampleGet{Name: &create.Name})
if err != nil {
	return nil, err
}

// Если есть
if e != nil {
	return nil, errors.New(`example already exist`)
}
```

Затем непосрдественно создание структуры и сохранение ее в БД

`ID` **всегда** должен быть в структуре и он **всегда** должен быть равен `uuid.NewString()`

```golang
// Создаем example
example := repositories.Example{
	ID:           uuid.NewString(),
	Name:         create.Name,
	DefaultValue: `default`,
}

res := r.db.Create(&example)
if res.Error != nil {
	return nil, res.Error
}
```

### GetOne

```golang
func (r *ExampleRepository) GetOne(query *repositories.ExampleGet) (*repositories.Example, error) {
	example := repositories.Example{}

	res := r.db.First(&example, query)

	// Если ничего не нашли
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	// Если ошибка
	if res.Error != nil {
		return nil, res.Error
	}

	return &example, nil
}
```

### GetAll

```golang
func (r *ExampleRepository) GetAll(query *repositories.ExampleGet) (*[]repositories.Example, error) {
	examples := []repositories.Example{}

	res := r.db.Find(&examples, query)

	// Если ничего не нашли
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	// Если ошибка
	if res.Error != nil {
		return nil, res.Error
	}

	return &examples, nil
}
```

### UpdateOne

```golang
func (r *ExampleRepository) UpdateOne(example *repositories.Example) error {
	res := r.db.Save(&example)
	return res.Error
}
```

### DeleteOne

```golang
func (r *ExampleRepository) DeleteOne(example *repositories.Example) error {
	res := r.db.Delete(&example, example)
	return res.Error
}
```

3. Инициализация репозитория

Далее в `server.go` в функцию `initReps` нужно добавить инициализацию только что созданного репозитория

```golang
example_rep.Init()
```

## Пример

Пример можно найти в `docs/examples/database`

- Папка `model` имитирует папку `repositories`
- Папка `example_rep` имитирует папку `example_rep`

В файле `example_rep.go` достаточно заменить слова `example` и `Example` на название модели, чтобы репозиторий начал быть рабочим

Импорт `example_model` **необходимо** удалить!