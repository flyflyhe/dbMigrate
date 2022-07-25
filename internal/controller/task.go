package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type task struct {
	base
}

func init() {
	t := &task{}
	SetRoute(&routeConfig{Method: http.MethodGet, Path: "/tasks", handle: t.List})
	SetRoute(&routeConfig{Method: http.MethodPost, Path: "/tasks", handle: t.Create})
}

func (this *task) List(c *gin.Context) {
	this.Success("hi list!", c)
}

func (this *task) Create(c *gin.Context) {
	this.Success("hi!", c)
}
