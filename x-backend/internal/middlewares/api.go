package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/jourloy/X-Backend/internal/repositories"
	"github.com/jourloy/X-Backend/internal/storage"
)

// Достает API ключ из заголовков и если есть - получает пользователя
func API() gin.HandlerFunc {
	return func(c *gin.Context) {
		api := c.Request.Header.Get(`api-key`)
		db := *storage.Database

		if api != `` {
			account := repositories.Account{}
			if res := db.First(&account, repositories.AccountGet{ApiKey: &api}); res.Error != nil {
				c.Next()
				return
			}

			if (account.ID != repositories.Account{}.ID) {
				c.Set(`account`, account)
				c.Set(`accountID`, account.ID)
				c.Set(`accountUsername`, account.Username)
			}
		}

		c.Next()
	}
}
