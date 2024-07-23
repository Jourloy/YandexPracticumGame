package guards

import (
	"github.com/gin-gonic/gin"
	"github.com/jourloy/X-Backend/internal/repositories"
)

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		a, exist := c.Get(`account`)

		// Если аккаунта нет
		if !exist {
			c.JSON(400, `api key is required or account is not exist`)
			c.Abort()
			return
		}

		account, ok := a.(repositories.Account)

		// Если аккаунт не валидный
		if !ok {
			c.JSON(400, `account is invalid`)
			c.Abort()
			return
		}

		// Если пользователь не админ
		if !account.IsAdmin {
			c.JSON(400, `only admin access`)
			c.Abort()
			return
		}

		c.Next()
	}
}
