package main

import (
	"corntab/master/config"
	"corntab/master/server"
	"github.com/sirupsen/logrus"
	"runtime"
)

func initEnv(){
	//配置线程数量
	runtime.GOMAXPROCS(runtime.NumCPU()-1)
}

func main(){
	//初始化线程
	initEnv()
	//首先加载配置文件
	config.Init()
	//初始化server
	server.Init()
	//启动Api HTTP服务
	if err := server.Run();err != nil{
		logrus.Errorf("run error")
	}
}