package response

import "net/http"

type APIResponse struct {
	Code  int         `json:"-"`
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

type ResponseBuilder struct {
	response *APIResponse
}

func NewResponse() *ResponseBuilder {
	return &ResponseBuilder{
		response: &APIResponse{},
	}
}

func (rb *ResponseBuilder) Code(code int) *ResponseBuilder {
	rb.response.Code = code
	return rb
}

func (rb *ResponseBuilder) Success(data interface{}) *ResponseBuilder {
	rb.response.Data = data
	return rb
}

func (rb *ResponseBuilder) Error(err *APIError) *ResponseBuilder {
	rb.response.Error = err
	return rb
}

func (rb *ResponseBuilder) Build() *APIResponse {
	return rb.response
}

func OK(data interface{}) *APIResponse {
	return NewResponse().Code(http.StatusOK).Success(data).Build()
}

func Created(data interface{}) *APIResponse {
	return NewResponse().Code(http.StatusCreated).Success(data).Build()
}

func NoContent() *APIResponse {
	return NewResponse().Code(http.StatusNoContent).Build()
}

func Err(code int, err *APIError) *APIResponse {
	return NewResponse().Code(code).Error(err).Build()
}
