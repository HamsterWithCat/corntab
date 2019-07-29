package main

import (
	"context"
	"corntab/prepare/etcd_connection"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"strconv"
	"time"
)

func main(){
	//放一个key-value键值对
	kv := clientv3.NewKV(etcd_connection.Client)
	putResp,err := kv.Put(context.TODO(),"/watch/test/hahaha","hahaha")
	if err != nil {
		fmt.Println("put key failed,reason:",err.Error())
		return
	}
	revisionId := putResp.Header.Revision // 操作的版本号

	go watherKey(context.TODO(),"/watch/test/hahaha",revisionId)

	go changeKey(context.TODO(),"/watch/test/hahaha")

	//主线程等待
	select {

	}

}
func watherKey(ctx context.Context,key string,revision_id int64){
	//获取watch
	watcher := clientv3.NewWatcher(etcd_connection.Client)
	//监听key 返回管道 从指定版本开始监听
	watchChan := watcher.Watch(ctx,key/*,clientv3.WithRev(revision_id)*/)
	//for循环监听管道的值 每一次变化被封装成一次WatchResponse响应放入管道内
	//但是 这里的如果使用了WithRev且设置的开始监听的变化不是下一次，那么第一次会返回多条已经存在的值 默认从下一次开始监听 值代表的本次也会被监听
	for watchResp :=range watchChan{
		fmt.Println("管道里的一条数据:",watchResp.CompactRevision,"事件数量:",len(watchResp.Events))
		for _,event := range watchResp.Events{
			fmt.Println("一次事件")
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("PUT:",event.Type,"  ",string(event.Kv.Key),"  ",string(event.Kv.Value))
			case mvccpb.DELETE:
				fmt.Println("DELETE:",event.Type,"  ",string(event.Kv.Key),"  ",string(event.Kv.Value))
			}
		}
		fmt.Println("--------------------------------------")
	}
}
func changeKey(ctx context.Context,key string){
	kv := clientv3.NewKV(etcd_connection.Client)
	curValue := 9
	for{
		curValue++
		kv.Put(ctx,key,"hahaha"+strconv.Itoa(curValue))
		time.Sleep(time.Second*1)
		kv.Delete(ctx,key)
		time.Sleep(time.Second*3)
	}
}