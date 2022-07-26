package controller

import (
	"github.com/flyflyhe/dbMigrate/internal/db"
	"github.com/flyflyhe/dbMigrate/internal/model"
	"github.com/flyflyhe/dbMigrate/internal/model/taskModel"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

type task struct {
	Base
}

func init() {
	t := &task{}
	SetRoute(&routeConfig{Method: http.MethodGet, Path: "/tasks", handle: t.List})
	SetRoute(&routeConfig{Method: http.MethodPost, Path: "/tasks", handle: t.Create})
	SetRoute(&routeConfig{Method: http.MethodPost, Path: "/task/:id/start", handle: t.Start})
}

func (this *task) List(c *gin.Context) {
	tModel := &taskModel.Task{}
	conn, err := tModel.GetDb()
	if err != nil {
		this.Failed("数据链接异常", c)
		log.Error().Caller().Err(err).Send()
		return
	}
	var tModelList []*taskModel.Task
	if err := conn.Table(tModel.TableName()).Find(&tModelList).Error; err != nil {
		this.Failed("数据链接异常", c)
		log.Error().Caller().Err(err).Send()
	}
	this.Success(tModelList, c)
}

func (this *task) Create(c *gin.Context) {
	var taskConfig *db.TaskExternalConfig
	var err error
	if err = c.ShouldBindJSON(&taskConfig); err != nil {
		this.Failed(err.Error(), c)
		log.Error().Caller().Err(err).Send()
		return
	}

	tModel := taskModel.CreateTaskByConfig(taskConfig)

	if err = model.Create(tModel); err != nil {
		this.Failed(err.Error(), c)
		log.Error().Caller().Err(err).Send()
		return
	}
	this.Success(tModel, c)
}

func (this *task) Start(c *gin.Context) {
	if c.Param("id") == "" {
		this.Failed("参数不能为空", c)
		return
	}

	if id, err := strconv.Atoi(c.Param("id")); err != nil {
		this.Failed(err.Error(), c)
	} else {
		tModel := &taskModel.Task{}
		if err := model.FindOne(tModel, int64(id)); err != nil {
			this.Failed(err.Error(), c)
			return
		}
		this.Success(tModel, c)
	}
}
