/**
 * @Author: Hardews
 * @Date: 2023/3/27 12:32
 * @Description:
**/

package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"spirit-core/model"
	"spirit-core/service"
	"spirit-core/tool"
)

func AddLetter(ctx *gin.Context) {
	var req struct {
		Content string `json:"content,omitempty"`
		Address string `json:"address,omitempty"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		tool.ParamFailed(ctx)
		return
	}

	if len(req.Content) > 2048 {
		tool.ErrorWithData(ctx, "文章过长")
		return
	}

	err = service.AddLetter(req.Address, req.Content)
	if err != nil {
		log.Println("add letter failed, err:", err)
		tool.InternetError(ctx)
		return
	}

	tool.OK(ctx)
}

func AddFutureLetter(ctx *gin.Context) {
	var fl model.LetterFutureForm
	err := ctx.ShouldBindJSON(&fl)
	if err != nil {
		tool.ParamFailed(ctx)
		log.Println(err)
		return
	}

	if len(fl.Content) > 2048 {
		tool.ErrorWithData(ctx, "信件过长")
		return
	}

	audit, err := service.AddFutureLetter(fl)
	if err != nil {
		log.Println("add future letter failed, err:", err)
		tool.InternetError(ctx)
		return
	}

	tool.OKWithData(ctx, audit)
}

func GetLetterSum(ctx *gin.Context) {
	tool.OKWithData(ctx, service.GetLetterSum())
}

func GetFutureLetter(ctx *gin.Context) {
	tool.OKWithData(ctx, service.GetFutureLetter())
}

func GetFutureLetterById(ctx *gin.Context) {
	tool.OKWithData(ctx, service.GetFutureLetterById(ctx.Query("id")))
}

func GetRandomFutureLetter(ctx *gin.Context) {
	tool.OKWithData(ctx, service.GetFutureRandomLetter())
}

func GetFutureLetterSum(ctx *gin.Context) {
	tool.OKWithData(ctx, service.GetFutureLetterSum())
}
