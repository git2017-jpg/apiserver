package response

import (
	"monitor-apiserver/pkg/errors"
	"monitor-apiserver/pkg/errors/code"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiResponse 代表一个响应给客户端的消息结构，包括错误码，错误消息，响应数据
type ApiResponse struct {
	RequestId string      `json:"request_id"`     // 请求的唯一ID
	ErrCode   int         `json:"err_code"`       // 错误码，0表示无错误
	Message   string      `json:"message"`        // 提示信息
	Data      interface{} `json:"data,omitempty"` // 响应数据，一般从这里前端从这个里面取出数据展示
}

//JSON 发送json格式的数据
func JSON(c *gin.Context, err error, data interface{}) {
	errCode, message := errors.DecodeErr(err)
	// 如果code != 0, 失败的话 返回http状态码400（一般也可以全部返回200）
	// 返回400 更严谨一些，个人接触的项目中大部分都是400。
	var httpStatus int
	if errCode != code.Success {
		httpStatus = http.StatusBadRequest
	} else {
		httpStatus = http.StatusOK
	}
	c.JSON(httpStatus, ApiResponse{
		RequestId: c.GetString("request_id"),
		ErrCode:   errCode,
		Message:   message,
		Data:      data,
	})
}
