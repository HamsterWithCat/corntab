package model

import "corntab/common"

type SaveJobReq struct {
	common.Job
}

type SaveJobResp struct {
	Job *common.Job `json:"job"`
}
