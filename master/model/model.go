package model

import "corntab/common"

//保存任务
type SaveJobReq struct {
	common.Job
}

type SaveJobResp struct {
	Job *common.Job `json:"job"`
}

//删除任务
type DeleteJobReq struct {
	JobName string `json:"job_name" form:"job_name" bind:"required"`
}

type DeleteJobResp struct {
	Job *common.Job `json:"job"`
}

//查询任务
type QueryJobReq struct {
	JobName string `json:"job_name" form:"job_name"`
}
type QueryJobResp struct {
	Jobs       []common.Job `json:"jobs"`
	TotalCount int          `json:"total_count"`
}

//强杀任务
type KillJobReq struct {
	JobName string   `json:"job_name"`
}

type KillJobResp struct {}