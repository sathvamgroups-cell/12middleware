package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("Request Method:", c.Request.Method)
		fmt.Println("Request Path:", c.Request.URL.Path)

		c.Next()
	}
}
