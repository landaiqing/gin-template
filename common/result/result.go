package result

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

const (
	SUCCESS = 0
	FAIL    = -1
)

func Result(code int, msg string, data any, success bool, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: msg,
		Data:    data,
		Success: success,
	})
}
func OK(msg string, data any, c *gin.Context) {
	Result(SUCCESS, msg, data, true, c)
}
func OkWithData(data any, c *gin.Context) {
	Result(SUCCESS, "success", data, true, c)
}
func OkWithMessage(msg string, c *gin.Context) {
	Result(SUCCESS, msg, nil, true, c)
}

func Fail(msg string, data any, c *gin.Context) {
	Result(FAIL, msg, data, false, c)
}
func FailWithMessage(msg string, c *gin.Context) {
	Result(FAIL, msg, nil, false, c)
}
func FailWithData(data any, c *gin.Context) {
	Result(FAIL, "fail", data, false, c)

}
func FailWithCode(code ErrorCode, c *gin.Context) {
	msg, ok := ErrMap[code]
	if ok {
		Result(int(code), msg, nil, false, c)
		return
	}
	Result(FAIL, "未匹配到对应的error code", nil, false, c)
}
