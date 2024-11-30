package errorx

const defaultCode = 1001

type CodeError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type CodeErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Message: msg}
}

func NewDefaultError(msg string) error {
	return NewCodeError(defaultCode, msg)
}

func (e *CodeError) Error() string {
	return e.Message
}

func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code:    e.Code,
		Message: e.Message,
	}
}
