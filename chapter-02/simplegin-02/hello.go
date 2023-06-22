package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func router() *gin.Engine {
	r := gin.Default()
	r.GET("/:name", func(c *gin.Context) {
		user := c.Param("name")
		c.String(200, fmt.Sprintf("hello, %s", user))
	})
	return r
}
func main() {
	router().Run()
}
