package controllers

import "net/http"

const (
	StatusParamError    = 10904001 //请求参数错误
	StatusDatabaseError = 10904002 // 数据库异常
	ServerError         = 10909999 //服务端异常
)

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(data interface{}) *Result {
	return &Result{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
	}
}

func FailWithData(code int, data interface{}) *Result {
	return &Result{
		Code:    code,
		Message: "FAIL",
		Data:    data,
	}
}

func FailWithMsg(code int, message string) *Result {
	return &Result{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}
