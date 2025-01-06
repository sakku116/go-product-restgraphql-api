package error_utils

type CustomErr struct {
	HttpCode int
	Message  string
	Detail   interface{}
	Data     interface{}
}

func (slf *CustomErr) Error() string {
	if slf.Detail != "" {
		err, ok := slf.Detail.(error)
		if ok {
			return err.Error()
		}

		detail, ok := slf.Detail.(string)
		if ok {
			return detail
		}

		return "unknown error"
	}
	return slf.Message
}
