package errors

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HttpError) Error() string {
	return e.Message
}
func NewCustomError(code int, message string) *HttpError {
	return &HttpError{
		Code:    code,
		Message: message,
	}
}

func ParseCustomError(err error) (int, string) {
	if customErr, ok := err.(*HttpError); ok {
		return customErr.Code, customErr.Message
	}
	return 500, err.Error()
}
