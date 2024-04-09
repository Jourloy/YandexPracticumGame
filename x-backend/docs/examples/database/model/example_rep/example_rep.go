package example_rep

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"gorm.io/gorm"

	example_model "github.com/jourloy/X-Backend/docs/examples/database/model"
	"github.com/jourloy/X-Backend/internal/storage"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[database-example]`,
		Level:  log.DebugLevel,
	})
)

var Repository example_model.ExampleRepository

type ExampleRepository struct {
	db gorm.DB
}

// Init создает репозиторий примера
func Init() {
	go migration()

	Repository = &ExampleRepository{
		db: *storage.Database,
	}
}

func migration() {
	if err := storage.Database.AutoMigrate(
		&example_model.Example{},
	); err != nil {
		logger.Fatal(`Migration failed`)
	}
}

func (r *ExampleRepository) Create(create *example_model.ExampleCreate) (*example_model.Example, error) {
	// Проверка, есть ли уже такой example
	e, err := r.GetOne(&example_model.ExampleGet{Name: &create.Name})
	if err != nil {
		return nil, err
	}

	// Если есть
	if e != nil {
		return nil, errors.New(`example already exist`)
	}

	// Создаем example
	example := example_model.Example{
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

func (r *ExampleRepository) GetOne(query *example_model.ExampleGet) (*example_model.Example, error) {
	example := example_model.Example{}

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

func (r *ExampleRepository) GetAll(query *example_model.ExampleGet) (*[]example_model.Example, error) {
	examples := []example_model.Example{}

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

func (r *ExampleRepository) UpdateOne(example *example_model.Example) error {
	res := r.db.Save(&example)
	return res.Error
}

func (r *ExampleRepository) DeleteOne(example *example_model.Example) error {
	res := r.db.Delete(&example, example)
	return res.Error
}
