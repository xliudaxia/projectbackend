package app

import (
	"bubble/pkg/errcode"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

type Error struct {
	code    int
	msg     string
	details []string
}

type Pager struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalRows int64 `json:"total_rows"`
}

func (r *Response) ToResponse(data interface{}) {
	hostname, _ := os.Hostname()
	if data == nil {
		data = gin.H{
			"code":      0,
			"msg":       "success",
			"tracehost": hostname,
		}
	} else {
		data = gin.H{
			"code":      0,
			"msg":       "success",
			"data":      data,
			"tracehost": hostname,
		}
	}
	r.Ctx.JSON(http.StatusOK, data)
}

func (r *Response) ToErrorResponse(err *errcode.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	r.Ctx.JSON(err.StatusCode(), response)
}

func (r *Response) ToResponseList(list interface{}, totalRows int64) {
	r.ToResponse(gin.H{
		"list": list,
		"pager": Pager{
			Page:      GetPage(r.Ctx),
			PageSize:  GetPageSize(r.Ctx),
			TotalRows: totalRows,
		},
	})
}
