package errorX

type Error struct {
	Code  ErrorCode `json:"code"`
	Error error     `json:"detail"`
}

func New(code ErrorCode) Error {
	return Error{
		Code:  code,
		Error: ERROR_MAP[code],
	}
}

func (ce Error) IsNotEmpty() bool {
	return ce != Error{Code: 0, Error: nil}
}

func (ce Error) GetErrorCode() int {
	return int(ce.Code)
}

func (ce Error) GetErrorCodeMessage() error {
	return ERROR_MAP[ce.Code]
}

func (ce Error) GetHttpCode() int {
	return ERROR_HTTP_MAP[ce.Code]
}
