package main

import (
	"github.com/gin-gonic/gin"
	"gocrazy/chapter-02/final-02/drawing"
)

func router() *gin.Engine {
	r := gin.Default()
	userRoute := r.Group("/image")
	{
		userRoute.GET("/:generator", func(c *gin.Context) {
			file := drawing.DrawOne("circles")
			c.Header("Content-Type", "image/png")
			c.File(file)
		})
	}
	return r
}

func main() {
	router().Run()
}
