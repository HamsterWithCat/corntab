package etcd_connection

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

var (
	config clientv3.Config
	Client *clientv3.Client
)

func init() {
	//读取客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"49.235.1.29:2379"},
		DialTimeout: time.Second * 5,
	}
	var err error
	if Client, err = clientv3.New(config); err != nil {
		panic(err)
	}
}
func main() {
	if Client != nil {
		fmt.Println("connect success")
	}
}
