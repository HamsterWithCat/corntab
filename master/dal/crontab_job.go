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
func (*JobMgr) SaveJob(jobName string, job string) (oldJob string, err error) {
	//获取操作etcd 客户端
	mgr := GetManager()
	kv := clientv3.NewKV(mgr.Client)

	var putResp *clientv3.PutResponse
	putResp, err = kv.Put(context.TODO(), jobName, job, clientv3.WithPrevKV())
	if err != nil {
		logrus.Errorf("put kv to etcd failed:err_msg = %v", err)
		return "", err
	}

	if putResp.PrevKv != nil {
		//获取修改之前的任务
		oldJob = string(putResp.PrevKv.Value)
	}

	return
}
