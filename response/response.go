package response

const (
	SuccessCode = "200"
	SuccessDesc = "Success."

	ErrBadRequestCode = "400"
	ErrBadRequestDesc = "Cannot bind request."

	ErrInternalServerCode = "500"
	ErrInternalServerDesc = "Internal service error."
)

type Responser interface {
	GetResponse() interface{}
}

type Response struct {
	Code        string      `json:"code" example:"200"`
	Description string      `json:"message" example:"Success"`
	Data        interface{} `json:"data,omitempty"`
}

type ErrResponse struct {
	Code        string      `json:"code" example:"400"`
	Description string      `json:"message" example:"Bad Request"`
	Error       interface{} `json:"error,omitempty"`
}

func NewResponse(code, desc string, data interface{}) *Response {
	return &Response{
		Code:        code,
		Description: desc,
		Data:        data,
	}
}

func NewErrResponse(code, desc string, err interface{}) *ErrResponse {
	return &ErrResponse{
		Code:        code,
		Description: desc,
		Error:       err,
	}
}

func (r *Response) GetResponse() interface{} {
	return r
}

func (e *ErrResponse) GetResponse() interface{} {
	return e
}
