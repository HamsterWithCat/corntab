package main

import (
	"corntab/master/config"
	"corntab/master/dal"
	"corntab/master/server"
	"github.com/sirupsen/logrus"
	"runtime"
)

func initEnv() {
	//配置线程数量
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
}

func main() {
	//初始化线程
	initEnv()
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message"},
	})
	//首先加载配置文件
	config.Init()
	//加载dal层配置
	dal.Init()
	//初始化server
	server.Init()
	//启动Api HTTP服务
	if err := server.Run(); err != nil {
		logrus.Errorf("run error")
		panic(err)
	}
}
