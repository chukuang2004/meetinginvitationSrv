package http

import (
	"beastpark/meetinginvitationservice/pkg/log"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// 通用的1开头
	Success            = 10000
	UnknownError       = 10001
	InvalidParameter   = 10002
	DBError            = 10003
	ServiceUnavailable = 10004
	DataNotFound       = 10005
)

var Code2Msg = map[int]string{
	Success:            "success",
	UnknownError:       "unknown error",
	InvalidParameter:   "invalid parameter",
	DBError:            "database error",
	ServiceUnavailable: "service unavailable",
	DataNotFound:       "data not found",
}

type Config struct {
	Port     int
	DNS      string
	SavePath string
}

func AppendErrorCode(m map[int]string) {
	for code, msg := range m {
		Code2Msg[code] = msg
	}
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WithResponse(ctx *gin.Context, code int, data interface{}, msg ...string) {
	s := Code2Msg[code]
	if len(msg) != 0 {
		s = strings.Join(msg, ";")
	}
	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: s,
		Data:    data,
	})
}

type ResponseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

func (w ResponseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

func LogHandler(ctx *gin.Context) {

	body := ""
	if !strings.Contains(ctx.GetHeader("content-type"), "multipart/form-data") {
		b, _ := ctx.GetRawData()
		body = string(b)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(b))
	}

	writer := ResponseWriter{
		ctx.Writer,
		bytes.NewBuffer([]byte{}),
	}
	ctx.Writer = writer

	ctx.Next()

	var str bytes.Buffer
	_ = json.Indent(&str, writer.b.Bytes(), "", "    ")

	log.GetInstance().Infof("Request %s, body %s, response status %d, body %s\n", ctx.Request.URL, body, ctx.Writer.Status(), str.String())

}
