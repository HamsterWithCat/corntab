package common

type Job struct {
	JobName  string `json:"job_name" form:"job_name" binding:"required"`
	Command  string `json:"command" form:"command" binding:"required"`
	CronExpr string `json:"cron_expr" form:"cron_expr" binding:"required"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponse() (resp Response) {
	resp = Response{
		Code: 0,
		Msg:  "",
		Data: struct {}{},
	}
	return
}
