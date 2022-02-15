package errors

type MyError struct {
	Code int
	Msg  string
	Data interface{}
}

var (
	LOGIN_UNKNOWN = NewError(202, "用户不存在")
	LOGIN_ERROR   = NewError(203, "账号或密码错误")
	VALID_ERROR   = NewError(300, "参数错误")
	ERROR         = NewError(400, "操作失败")
	UNAUTHORIZED  = NewError(401, "您还未登录")
	NOT_FOUND     = NewError(404, "资源不存在")
	INNER_ERROR   = NewError(500, "系统发生异常")
)

func (e *MyError) Error() string {
	return e.Msg
}

func NewError(code int, msg string) *MyError {
	return &MyError{
		Msg:  msg,
		Code: code,
	}
}
