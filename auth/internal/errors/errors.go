package errors

type HttpError struct {
	StatusCode int
	Message    string
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewHttpError(status int, msg string) *HttpError {
	return &HttpError{
		StatusCode: status,
		Message:    msg,
	}
}

func ParseHttpError(err error) (int, string, error) {
	if httpErr, ok := err.(*HttpError); ok {
		return httpErr.StatusCode, httpErr.Message, nil
	}
	return 0, "", err
}
