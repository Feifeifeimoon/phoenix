package common

// 统一返回的json字段
type Response struct {
	Success int         `json:"success"`
	Data    interface{} `json:"data"`
	Count   int         `json:"count"`
	Msg     string      `json:"msg"`
}

func SuccessResponse(data interface{}, count int) *Response {
	return &Response{
		Success: 1,
		Data:    data,
		Count:   count,
		Msg:     "",
	}
}

func ErrorResponse(msg string) *Response {
	return &Response{
		Success: 0,
		Data:    nil,
		Count:   0,
		Msg:     msg,
	}
}
