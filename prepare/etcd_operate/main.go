package main

import (
	"context"
	"corntab/prepare/etcd_connection"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

func main() {
	getKeyValue()
	//deleteKey()
}
func putKeyValue() {
	//获取KV,用于读写etcd的键值对
	kv := clientv3.NewKV(etcd_connection.Client)
	//写操作
	resp, err := kv.Put(context.TODO(), "/cron/jobs/job1", "hello")
	if err != nil {
		fmt.Printf("put key failed,reason[%v]", err.Error())
		return
	}
	//查看版本号
	fmt.Println("Revision:", resp.Header.Revision)
	//验证key的多版本控制
	//WithPrevKV 获取之前版本的key value
	resp, err = kv.Put(context.TODO(), "/cron/jobs/job1", "hi", clientv3.WithPrevKV())
	if err != nil {
		fmt.Printf("put key[%v] failed,reason[%v]", "/cron/jobs/job1", err.Error())
		return
	}
	fmt.Println("Revision:", resp.Header.Revision)
	fmt.Println("PreKV:", string(resp.PrevKv.Key), " ", string(resp.PrevKv.Value), " ", resp.PrevKv.Version)
}
func getKeyValue() {
	//WithCountOnly //返回key的个数
	//WithPrefix	//指定操作为前缀操作
	//With
	//WithLimit	    //限定搜索个数
	//WithFromKey   //从指定的key开始获取
	kv := clientv3.NewKV(etcd_connection.Client)

	resp, err := kv.Get(context.TODO(), "/cron/jobs/")
	if err != nil {
		fmt.Printf("get key failed ,reason[%v]", err.Error())
		return
	}
	//获取到key值
	//cron/jobs/job1 hi
	//key不存在 通过resp.Count<1判断
	if resp.Count > 0 {
		for _, pairkv := range resp.Kvs {
			fmt.Println(string(pairkv.Key) + " " + string(pairkv.Value))
		}
	} else {
		//key不存在的情况
		fmt.Println("server doesn't has this key's record")
	}

	//获取指定前缀的所有key-value
	resp, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("get key with prefix failed,reason[%v]\n", err.Error())
		return
	}
	if resp.Count > 0 {
		for _, kvpair := range resp.Kvs {
			fmt.Println(kvpair.Key, "  ", kvpair.Value)
		}
	} else {
		fmt.Println("server doesn't has this with prefix key's record")
	}
}
func deleteKey() {
	kv := clientv3.NewKV(etcd_connection.Client)
	//删除指定key
	//DeleteResponse具有Header(etcd连接相关)、Deleted(影响key数)、PrevKvs(影响的key的前版本)
	resp, err := kv.Delete(context.TODO(), "/cron/jobs/", clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("delete key failed,reason[%v]\n", err.Error())
		return
	}
	if resp.Deleted > 0 {
		fmt.Printf("delete %d keys\n", resp.Deleted)
	} else {
		fmt.Println("no key deleted")
	}

}
