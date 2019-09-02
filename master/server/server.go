package server

import (
	"corntab/master/config"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
)

var (
	server *gin.Engine
)

//初始化server
func Init() {
	//创建路由
	server = gin.New()
	//定义映射
	server.POST("/task/save", saveTask)
	server.GET("/ping", pong)
	server.POST("/task/delete", deleteTask)
	server.GET("/task/query", queryTask)
	server.POST("/task/kill",killTask)
}
func Run() error {
	if server == nil {
		logrus.Errorf("[Server.Run] server not initialized")
		return errors.New("server not initialized")
	}
	err := server.Run(":" + strconv.Itoa(int(config.GetServerConfig().IP)))
	if err != nil {
		logrus.Errorf("[Server.Run]run server failed:%v", err)
		return err
	}

	logrus.Infof("[Server.Run]server is running")
	return nil
}
