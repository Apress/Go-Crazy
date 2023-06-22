package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
	"gocrazy/chapter-02/final-03/drawing"
	"golang.org/x/exp/maps"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type jobData struct {
	Id        string
	Generator string
}

var sm sync.Map

func (j *jobData) Bytes() []byte {
	b, _ := json.Marshal(j)
	return b
}

func router() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")

	rand.Seed(time.Now().Unix())

	q := queue.NewPool(30, queue.WithFn(func(ctx context.Context, m core.QueuedMessage) error {
		j, _ := m.(*jobData)
		json.Unmarshal(m.Bytes(), &j)

		sleepTime := time.Duration(rand.Intn(10)) * time.Second
		time.Sleep(sleepTime)
		path := drawing.DrawOne(j.Generator)
		sm.Store(j.Id, path)
		fmt.Printf("Stored: %s:%s [%s]\n", j.Id, j.Generator, path)

		return nil
	}))

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
	newRoute := r.Group("/new")
	{
		newRoute.GET("/load/:id", func(c *gin.Context) {
			id := c.Param("id")
			path, ok := sm.Load(id)

			if ok {
				fmt.Printf("Found %s for id: %s\n", path, id)
				c.Header("Content-Type", "image/png")
				c.File(fmt.Sprintf("%s", path.(string)))
			} else {
				fmt.Printf("Path not found for id: %s\n", id)
				c.Header("Content-Type", "image/jpg")
				c.Header("Cache-Control", "no-cache")
				c.File("static/loading.jpg")
			}
		})
		newRoute.GET("/:generator", func(c *gin.Context) {
			generator := c.Param("generator")
			newJob := jobData{
				Id:        strconv.Itoa(rand.Int()),
				Generator: generator,
			}
			q.Queue(&newJob)
			res := map[string]string{"id": newJob.Id, "url": "http://" + c.Request.Host + "/new/load/" + newJob.Id}
			c.JSON(200, res)
		})
	}

	return r
}

func main() {
	router().Run()
}
