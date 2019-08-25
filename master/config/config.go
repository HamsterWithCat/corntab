package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"runtime"
	"strings"
)

var (
	cfg *serverConfig
)

func GetServerConfig()*serverConfig{
	if cfg == nil {
		panic("config not initialized")
	}
	return cfg
}

type serverConfig struct {
	IP     int32     `yaml:"server_ip"`
	//连接超时时间
	ConnectionTimeout int64 `yaml:"server_connection_timeout"`
}

//加载配置文件
func Init(){
	//加载出错 panic
	if err := parseConf();err != nil{
		panic(err)
	}
}

//解析配置文件
func parseConf()error{
	cfg = &serverConfig{}

	//配置文件名称
	confFileName := "corntab_server.yml"

	//拼接路径 根据当前配置文件的绝对路径解析出conf文件夹路径，config和conf两文件夹位置不能改变
	configPath := getCurFilePath()
	pos := strings.LastIndex(configPath,"/")
	rootFilePath := configPath[:pos]
	confFilePath := rootFilePath+"/conf/"+confFileName

	//读取配置文件信息
	confStream ,err := ioutil.ReadFile(confFilePath)
	if err != nil {
		logrus.Errorf("[parseConfig] read conf file failed:%v",err)
		return err
	}

	//解析配置文件
	err = yaml.Unmarshal(confStream,cfg)
	if err != nil {
		logrus.Errorf("[parseConfig] unmarshal config file failed:")
		return err
	}

	return nil
}

//获取当前文件夹路径
func getCurFilePath() string {
	_, fileLoc, _, ok := runtime.Caller(0)
	if !ok {
		logrus.Errorf("[Test] get file location failed!")
		return ""
	}
	pos := strings.LastIndex(fileLoc, "/")
	curFolder := fileLoc[0 : pos]
	return curFolder
}
