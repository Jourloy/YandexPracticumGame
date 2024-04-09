package guards

import (
	"github.com/gin-gonic/gin"
)

// Достает пользователя из контекста. Если нет - отказ
func CheckAPI() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exist := c.Get(`account`)

		// Если аккаунта нет
		if !exist {
			c.JSON(400, gin.H{`error`: `api key is required or account is not exist`})
			c.Abort()
			return
		}

		c.Next()
	}
}
