package Proto_Proxy

// 为proto中定义的错误实现Error接口
func (e *ErrorST) Error() string {
	return e.Msg
}
