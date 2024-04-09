package middlewares

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

var (
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: `[request]`,
		Level:  log.DebugLevel,
	})
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		logger.Info(
			`Response`,
			`method`, method,
			`path`, path,
			`status`, status,
			`latency`, latency,
		)

		// Если запрос длится более 300 миллисекунд
		if time.Duration(latency.Milliseconds()) > 300 {
			logger.Warn(`Latency is over 300ms`)
			// TODO: запись в БД или отправка webhook
		}
	}
}
