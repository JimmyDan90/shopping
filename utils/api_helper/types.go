package api_helper

import "errors"

// Response 响应结构体
type Response struct {
	Message string `json:"message"`
}

// ErrorResponse 响应错误结构体
type ErrorResponse struct {
	Message string `json:"errorMessage"`
}

// ErrInvalidBody 自定义错误
var (
	ErrInvalidBody = errors.New("请检查你的请求体")
)
