package controller

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"sync"
)

var routeDepository []*routeConfig
var routeLock sync.Mutex

type routeConfig struct {
	Method     string
	Path       string
	handle     func(c *gin.Context)
	middleware []gin.HandlerFunc
}

func SetRoute(config *routeConfig) {
	routeLock.Lock()
	defer routeLock.Unlock()
	routeDepository = append(routeDepository, config)
}

type base struct {
}

func (this *base) Success(data any, c *gin.Context) {
	c.JSON(200, data)
}

func (this *base) Failed(msg string, c *gin.Context) {
	c.JSON(500, gin.H{"err": msg})
}

func Start(addr string) {
	gin.DefaultWriter = io.MultiWriter(os.Stdout)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	for _, config := range routeDepository {
		router.Use(config.middleware...).Handle(config.Method, config.Path, config.handle)
	}

	router.Run(addr)
}
