package main

import (
	"context"
	"corntab/prepare/etcd_connection"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	//获取上下文信息
	ctxRoot := context.Background()

	//新建租约 并定时对租约进行续约
	lease := clientv3.NewLease(etcd_connection.Client)
	//5second 租期的租约
	leaseResp, err := lease.Grant(context.TODO(), 5)
	if err != nil {
		fmt.Println("generate lease failed,reason ", err.Error())
		return
	}
	//派生上下文，用于对续约的管理
	ctx, cancelFunc := context.WithCancel(ctxRoot)
	//函数退出时 关闭自动续约
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseResp.ID)
	//对租期自动续约
	go autoKeepAlive(ctx, leaseResp.ID)

	//开启事务
	kv := clientv3.NewKV(etcd_connection.Client)
	txn := kv.Txn(context.TODO())
	//如果不存在锁，则上锁
	key := "/cron/lock/hahalock"
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut(key, "I'm lock", clientv3.WithLease(leaseResp.ID))).
		Else(clientv3.OpGet(key))
	txnResp, err := txn.Commit()
	if err != nil {
		fmt.Println("transaction commit failed,err ", err.Error())
		return
	}

	//判断事务的执行结果
	if txnResp.Succeeded {
		//加锁成功，执行逻辑
		time.Sleep(time.Second * 5)
		fmt.Println("----------hahahahahahaha---------------")
	} else {
		fmt.Println("lock has been localed")
		//注意 使用txnResp获取值时，可以获取put、get、delete的执行结果，但前提是在事务中执行了这些语句，有但没执行也不可以
		fmt.Println(string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value)) //输出当前锁的值
	}
}

//使用此种方法可能会产生问题，在请求续约时可能会失败，但是进行的处理只是直接返回，不进行续约，而这对于业务主函数是不可见的
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
