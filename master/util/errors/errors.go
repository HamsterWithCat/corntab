package errors

//错误结构
const (
	SUCCESS = iota
	UNDEFINED
	PARAMBINDERROR
	PARAMPARSEERROR
	TASKNOTEXIST

	SAVETASKERROR
	DELETETASKERROR
	QUERYTASKERROR
	KILLTASKERROR

	CRONEXPRESSIONERROR
)

var (
	errMap = map[int]string{
		SUCCESS:         "",
		UNDEFINED:       "系统内部错误",
		PARAMBINDERROR:  "参数绑定错误",
		PARAMPARSEERROR: "参数解析失败",
		TASKNOTEXIST:    "任务不存在",

		SAVETASKERROR:   "存储任务失败,请重试",
		DELETETASKERROR: "删除任务失败,请重试",
		QUERYTASKERROR:  "查询任务失败,请重试",
		KILLTASKERROR:   "杀死任务失败，请重试",

		CRONEXPRESSIONERROR: "cron表达式解析失败",
	}
)

type CTErr struct {
	code int
}

func (err CTErr) Error() string {
	return errMap[err.code]
}

func NewCTErr(code int) CTErr {
	err := CTErr{
		code: code,
	}
	return err
}

func GetMsg(code int) string {
	msg, ok := errMap[code]
	if !ok {
		return errMap[UNDEFINED]
	}
	return msg
}

//获取错误信息
func GetErr(err error) (code int, msg string) {
	ctErr, ok := err.(CTErr)

	if !ok {
		return UNDEFINED, errMap[UNDEFINED]
	}
	return ctErr.code, errMap[ctErr.code]
}
