package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Message struct {
	// json tag to de-serialize json body
	Name  string `json:"name"`
	Email string `json:"email" binding:"required,email"`
}

func router() *gin.Engine {
	r := gin.Default()
	userRoute := r.Group("/user")
	{
		userRoute.GET("/hello/:name", func(c *gin.Context) {
			user := c.Param("name")
			response := fmt.Sprintf("hello, %s", user)
			c.String(http.StatusOK, response)
		})
		userRoute.POST("/post", func(c *gin.Context) {
			body := Message{}
			if err := c.BindJSON(&body); err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}
			fmt.Println(body)
			c.JSON(http.StatusAccepted, &body)
		})
		userRoute.POST("/upload", func(c *gin.Context) {
			file, _ := c.FormFile("file")
			log.Println(file.Filename)

			// Upload the file to specific dst.
			c.SaveUploadedFile(file, "/tmp/tempfile")

			c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
		})
	}
	return r
}

func main() {
	router().Run()
}
