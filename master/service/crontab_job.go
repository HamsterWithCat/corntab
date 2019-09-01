package service

import (
	"corntab/master/dal"
	"corntab/master/model"
	"corntab/master/util"
	"corntab/master/util/errors"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"github.com/sirupsen/logrus"
)

//参数为指针类型 返回值为值类型时 函数执行速度最快
func SaveJob(job *model.SaveJobReq) (oldJob model.SaveJobResp, err error) {
	logrus.Infof("save task to etcd server,job name : %s", job.JobName)

	oldJob = model.SaveJobResp{}
	//构造返回值
	//判断corn表表达式是否合规
	_, err = cronexpr.Parse(job.CronExpr)
	if err != nil {
		logrus.Warnf("[service.SaveJob] cron expression is invalid,err_msg = %v", err)
		return oldJob, errors.NewCTErr(errors.CRONEXPRESSIONERROR)
	}

	//任务加密为字符串
	bytes, err := json.Marshal(job.Job)
	if err != nil {
		logrus.Warnf("[service.SaveJob] marshal task failed,err_msg = %v", err)
		return oldJob, errors.NewCTErr(errors.PARAMPARSEERROR)
	}
	//存储任务
	jobName := util.CORNJOBNAMEPREFIX + job.JobName
	old, err := dal.GetJobMgr().SaveJob(jobName, string(bytes))
	if err != nil {
		logrus.Errorf("[service.SaveJob] save task to etcd serve failed,err_msg = %v", err)
		return oldJob, errors.NewCTErr(errors.SAVETASKERROR)
	}

	//解析返回值
	if "" != old {
		err = json.Unmarshal([]byte(old), &oldJob.Job)
		if err != nil {
			logrus.Warnf("[service.SaveJob] unmarshal task failed,err_msg = %v", err)
		}
	}
	//返回成功
	return oldJob, nil
}
