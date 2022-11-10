package utils

type Response struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   any    `json:"data"`
}
