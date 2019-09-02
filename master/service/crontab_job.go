package service

import (
	"corntab/common"
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
	//构造返回值
	oldJob = model.SaveJobResp{}
	//check name
	if "" == job.JobName {
		logrus.Errorf("[service.SaveJob] job name is empty,can not create")
		return oldJob, errors.NewCTErr(errors.PARAMPARSEERROR)
	}

	logrus.Infof("save task to etcd server,job name : %s", job.JobName)

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
	if nil != old {
		err = json.Unmarshal([]byte(old), &oldJob.Job)
		if err != nil {
			logrus.Warnf("[service.SaveJob] unmarshal task failed,err_msg = %v", err)
		}
	}
	//返回成功
	return oldJob, nil
}

func DeleteJob(job *model.DeleteJobReq) (deleteJob model.DeleteJobResp, err error) {
	deleteJob = model.DeleteJobResp{}
	//判断job name是否为空
	if "" == job.JobName {
		logrus.Errorf("[service.DeleteJob] job name is empty,not allow delete all task")
		return deleteJob, errors.NewCTErr(errors.TASKNOTEXIST)
	}

	logrus.Infof("[service.DeleteJob] delete job which key is %s", job.JobName)

	jobName := util.CORNJOBNAMEPREFIX + job.JobName
	old, err := dal.GetJobMgr().DeleteJob(jobName)
	if err != nil {
		logrus.Warnf("[service.DeleteJob]delete job failed,job_name = %s,err_msg = %v", jobName, err)
		return deleteJob, errors.NewCTErr(errors.DELETETASKERROR)
	}

	if nil == old {
		return deleteJob, errors.NewCTErr(errors.TASKNOTEXIST)
	}

	err = json.Unmarshal(old, &deleteJob.Job)
	if err != nil {
		logrus.Warnf("[service.DeleteTask] unmarshal task failed,job_name = %s,err_msg = %v", jobName, err)
	}
	return deleteJob, nil
}

func ListJobs(req *model.QueryJobReq) (queryJobs model.QueryJobResp, err error) {
	queryJobs = model.QueryJobResp{}
	logrus.Infof("[service.ListJobs] query job list with prefix:%s", req.JobName)
	//解析参数
	jobName := util.CORNJOBNAMEPREFIX + req.JobName

	jobs, err := dal.GetJobMgr().QueryJobWithPrefix(jobName)
	if err != nil {
		logrus.Errorf("[service.ListJobs]query job list error,prefix key = %v,err_msg = %v", req.JobName, err)
		return queryJobs, errors.NewCTErr(errors.QUERYTASKERROR)
	}

	for _, jobBytes := range jobs {
		job := common.Job{}
		err = json.Unmarshal(jobBytes, &job)
		if err != nil {
			logrus.Warnf("[service.ListJobs]unmarshal job failed,err_msg = %v", err)
			continue
		}
		queryJobs.Jobs = append(queryJobs.Jobs, job)
	}

	queryJobs.TotalCount = len(queryJobs.Jobs)
	return queryJobs, nil
}

//杀死任务
func KillJob(req *model.KillJobReq)(resp model.KillJobResp,err error){
	resp = model.KillJobResp{}
	logrus.Infof("[service.KillJob] kill job,job_name = %s",req.JobName)

	if req.JobName == ""{
		logrus.Warnf("[service.KillJob] can't kill job with empty job name")
		return resp,errors.NewCTErr(errors.TASKNOTEXIST)
	}
	jobName := util.CORNJOBKILLNAMEPREFIX+req.JobName
	err = dal.GetJobMgr().SaveJobNameWithLease(jobName)
	if err != nil {
		logrus.Warnf("[service.KillJob] kill job failed,err_msg = %v",err)
		return resp,errors.NewCTErr(errors.KILLTASKERROR)
	}
	return resp,nil
}