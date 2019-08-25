package main

import (
	"context"
	"corntab/prepare/etcd_connection"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	executeFunc(wg)

	wg.Wait()

	time.Sleep(time.Second * 5)
}
func executeFunc(wg *sync.WaitGroup) {
	defer wg.Done()
	//申请一个租约
	lease := clientv3.NewLease(etcd_connection.Client)

	//申请一个10秒的租约
	leaseResp, err := lease.Grant(context.TODO(), 10)
	if err != nil {
		fmt.Printf("apply grant failed,reason[%v]", err.Error())
		return
	}
	//使用租约的id创建一个key
	kv := clientv3.NewKV(etcd_connection.Client)

	_, err = kv.Put(context.TODO(), "/lease/test/hahaha", "hahaha", clientv3.WithLease(leaseResp.ID))
	if err != nil {
		fmt.Printf("put key with lease failed,reason[%s]\n", err.Error())
		return
	}
	//获取context
	ctxRoot := context.Background()
	ctx, cancelFunc := context.WithCancel(ctxRoot)
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseResp.ID)

	//执行续约操作
	go autoKeepAlive(ctx, leaseResp.ID)
	cancelFunc = cancelFunc //防止报错
	//轮训判断key是否过期
	for i := 0; i < 10; i++ {
		//查找key
		getResp, err := kv.Get(context.TODO(), "/lease/test/hahaha")
		//获取出错
		if err != nil {
			fmt.Printf("get key error,reason[%s]\n", err.Error())
			return
		}
		//查询成功
		if getResp.Count > 0 {
			for _, kvpair := range getResp.Kvs {
				fmt.Println(string(kvpair.Key), "  ", string(kvpair.Value))
			}
		} else {
			fmt.Println("key isn't exist")
			break
		}
		//两秒查询一次
		time.Sleep(time.Second * 2)
	}
	//退出时 先取消续约 后删除key
	cancelFunc()                               //停止续约
	lease.Revoke(context.TODO(), leaseResp.ID) //停止租约

	//kv.Delete(context.TODO(),"/lease/test/hahaha")
	fmt.Println("execute goroutine exit")
}
func autoKeepAlive(ctx context.Context, id clientv3.LeaseID) {
	//获取租期操作对象
	lease := clientv3.NewLease(etcd_connection.Client)
	//循环执行续约操作
	kaResp, err := lease.KeepAlive(ctx, id)
	if err != nil {
		fmt.Printf("keep alive for id[%d] failed", id)
		return
	}

	for reply := range kaResp {
		if reply == nil {
			return
		}
		fmt.Println(reply.ID)
	}
	fmt.Println("keep alive function exit")
}
