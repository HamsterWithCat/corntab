package errors

//错误结构
const (
	SUCCESS = iota
	UNDEFINED
	PARAMBINDERROR
	PARAMPARSEERROR

	SAVETASKERROR

	CRONEXPRESSIONERROR
)

var (
	errMap = map[int]string{
		SUCCESS:         "",
		UNDEFINED:       "系统内部错误",
		PARAMBINDERROR:  "参数绑定错误",
		PARAMPARSEERROR: "参数解析失败",

		SAVETASKERROR: "存储任务失败",

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
