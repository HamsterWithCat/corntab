package dal

import (
	"corntab/master/config"
	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	mgr  *manager
	Once sync.Once
)

type manager struct {
	Client *clientv3.Client
}

//dal层控制器 类似数据库连接管理
func Init() {
	serverCfg := config.GetServerConfig()

	//读取配置信息
	etcdConfig := clientv3.Config{
		Endpoints:   serverCfg.EtcdEndpoints,
		DialTimeout: time.Second * time.Duration(serverCfg.EtcdDialTimeout),
	}

	client, err := clientv3.New(etcdConfig)
	if err != nil {
		logrus.Errorf("init etcd client failed,err_msg = %v", err)
		panic(err)
	}
	mgr = &manager{}
	mgr.Client = client
}

func GetManager() *manager {
	return mgr
}
