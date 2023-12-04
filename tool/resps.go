/**
 * @Author: Hardews
 * @Date: 2023/3/19 22:11
 * @Description:
**/

package tool

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type respForm struct {
	Status int    `json:"status,omitempty"`
	Msg    string `json:"msg,omitempty"`
}

var (
	LoginNeed = respForm{
		Status: http.StatusOK,
		Msg:    "need login",
	}
	Forbidden = respForm{
		Status: http.StatusForbidden,
		Msg:    "无权限操作",
	}
	ParamError = respForm{
		Status: http.StatusOK,
		Msg:    "param error",
	}
	SqlError = respForm{
		Status: http.StatusInternalServerError,
		Msg:    "sql error",
	}
)

func OK(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "successful",
	})
}

func OKWithData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    data,
	})
}

func InternetError(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "internet error",
	})
}

func OperateForbid(ctx *gin.Context) {
	ctx.JSON(Forbidden.Status, Forbidden)
}

func ParamFailed(ctx *gin.Context) {
	ctx.JSON(ParamError.Status, ParamError)
}

func SQLFailed(ctx *gin.Context) {
	ctx.JSON(SqlError.Status, SqlError)
}

func NeedLogin(ctx *gin.Context) {
	ctx.JSON(LoginNeed.Status, LoginNeed)
}

func ErrorWithData(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"status": http.StatusInternalServerError,
		"msg":    data,
	})
}
