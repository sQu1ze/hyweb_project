package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Error     interface{} `json:"error"`
	Code      int         `json:"code"`
	Timestamp string      `json:"timestamp"`
}

func JSON(c *gin.Context, httpStatus int, success bool, message string, data interface{}, err interface{}) {
	c.JSON(httpStatus, Response{
		Success:   success,
		Message:   message,
		Data:      data,
		Error:     err,
		Code:      httpStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	JSON(c, http.StatusOK, true, message, data, nil)
}

func Error(c *gin.Context, httpStatus int, message string, err interface{}) {
	JSON(c, httpStatus, false, message, nil, err)
}
