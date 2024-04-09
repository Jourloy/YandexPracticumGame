package tools

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/gin-gonic/gin"
)

// Переводит body в структуру
func ParseBody(c *gin.Context, body interface{}) error {
	// Проверка body
	if c.Request.Body == nil {
		return errors.New(`body not found`)
	}

	// Чтение body
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	// Парсинг
	if err := json.Unmarshal(b, &body); err != nil {
		return err
	}

	return nil
}

// Contains проверяет наличие элемента "e" с типом "T" в массиве "s"
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
