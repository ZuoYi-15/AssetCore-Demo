package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "success", Data: data, TraceID: traceID(c)})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Body{Code: 0, Message: "success", Data: data, TraceID: traceID(c)})
}

func Fail(c *gin.Context, status int, code int, message string) {
	c.JSON(status, Body{Code: code, Message: message, TraceID: traceID(c)})
}

func traceID(c *gin.Context) string {
	v, _ := c.Get("trace_id")
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
