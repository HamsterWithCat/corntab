package dal

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/sirupsen/logrus"
)

type JobMgr struct{}

var (
	jobMgr *JobMgr
)

func GetJobMgr() *JobMgr {
	Once.Do(func() { jobMgr = &JobMgr{} })
	return jobMgr
}

//保存任务
func (*JobMgr) SaveJob(key string, value string) (oldValue []byte, err error) {
	//获取操作etcd 客户端
	mgr := GetManager()
	kv := clientv3.NewKV(mgr.Client)

	var putResp *clientv3.PutResponse
	putResp, err = kv.Put(context.TODO(), key, value, clientv3.WithPrevKV())
	if err != nil {
		logrus.Errorf("put kv to etcd failed:key = %s ,err_msg = %v", key, err)
		return nil, err
	}

	if putResp.PrevKv != nil {
		//获取修改之前的任务
		oldValue = putResp.PrevKv.Value
	}

	return
}

//删除任务
func (*JobMgr) DeleteJob(key string) (oldValue []byte, err error) {
	//获取操作etcd 客户端
	mgr := GetManager()
	kv := clientv3.NewKV(mgr.Client)

	var deleteResp *clientv3.DeleteResponse
	deleteResp, err = kv.Delete(context.TODO(), key, clientv3.WithPrevKV())
	if err != nil {
		logrus.Errorf("delete key failed:key = %s ,err_mag = %v", key, err)
		return nil, err
	}

	if len(deleteResp.PrevKvs) != 0 {
		oldValue = deleteResp.PrevKvs[0].Value
	}

	return
}

//查询所有任务
func (*JobMgr) QueryJobWithPrefix(key string) (values [][]byte, err error) {
	//获取操作etcd 客户端
	mgr := GetManager()
	kv := clientv3.NewKV(mgr.Client)

	var getResp *clientv3.GetResponse
	getResp, err = kv.Get(context.TODO(), key, clientv3.WithPrefix())
	if err != nil {
		logrus.Errorf("query job failed,job prefix is %s,err_msg = %v", key, err)
		return nil, err
	}

	values = make([][]byte, 0, len(getResp.Kvs))
	for _, result := range getResp.Kvs {
		values = append(values, result.Value)
	}

	return values, nil
}
