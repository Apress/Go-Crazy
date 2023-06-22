package main

import (
	"github.com/gin-gonic/gin"
	"gocrazy/chapter-02/final-02/drawing"
	"golang.org/x/exp/maps"
	"net/http"
)

func router() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")

	imageRoute := r.Group("/image")
	{
		imageRoute.GET("/:generator", func(c *gin.Context) {
			generator := c.Param("generator")
			file := drawing.DrawOne(generator)
			c.Header("Content-Type", "image/png")
			c.File(file)
		})
	}
	listRoute := r.Group("/list")
	{
		listRoute.GET("/simple", func(c *gin.Context) {
			c.HTML(http.StatusOK, "simple.tmpl", gin.H{
				"keys": maps.Keys(drawing.DRAWINGS),
			})
		})
		listRoute.GET("/bootstrap", func(c *gin.Context) {
			c.HTML(http.StatusOK, "bootstrap.tmpl", gin.H{
				"keys": maps.Keys(drawing.DRAWINGS),
			})
		})
	}

	return r
}

func main() {
	router().Run()
}
