package example_model

import (
	"time"

	"gorm.io/gorm"
)

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

type ExampleCreate struct {
	Name      string `json:"name"`
	ExampleID string `json:"exmapleId"`
	ParentID  string `json:"parentId"`
}

type ExampleGet struct {
	ID        *string
	Name      *string
	ExampleID *string
	ParentID  *string
	Limit     *int
}

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

type Child struct{}
