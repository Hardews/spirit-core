/**
 * @Author: Hardews
 * @Date: 2023/3/25 14:35
 * @Description:
**/

package api

import (
	"errors"
	"spirit-core/my_consts"
	"spirit-core/service"
	"spirit-core/tool"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"log"
)

func GetReplyLetterAdvice(ctx *gin.Context) {
	admId, res := ctx.Get("gmt_id")
	if !res {
		tool.OperateForbid(ctx)
		return
	}

	replyId, res := ctx.GetQuery("id")
	if !res {
		tool.ParamFailed(ctx)
		return
	}

	letter, err := service.GetMyReplyInfo(replyId, admId.(string))
	if err != nil {
		if errors.Is(err, my_consts.ErrOfForbidden) {
			tool.OperateForbid(ctx)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tool.OKWithData(ctx, "无此回复信息")
			return
		}
		tool.InternetError(ctx)
		log.Println("get letter info failed,err:", err)
		return
	}

	tool.OKWithData(ctx, letter)
}

func GetMyReplyOverview(ctx *gin.Context) {
	admId, res := ctx.Get("gmt_id")
	if !res {
		tool.OperateForbid(ctx)
		return
	}

	letters, err := service.GetMyReplyLetter(admId.(string))
	if err != nil {
		if errors.Is(err, my_consts.ErrOfForbidden) {
			tool.OperateForbid(ctx)
			return
		}
		tool.InternetError(ctx)
		log.Println("get letter reply homepage failed,err:", err)
		return
	}

	tool.OKWithData(ctx, letters)
}

func GetLetterReplyHomepage(ctx *gin.Context) {
	admId, res := ctx.Get("gmt_id")
	if !res {
		tool.OperateForbid(ctx)
		return
	}

	letters, err := service.GetReplyLetterHomepage(admId.(string))
	if err != nil {
		if errors.Is(err, my_consts.ErrOfForbidden) {
			tool.OperateForbid(ctx)
			return
		}
		tool.InternetError(ctx)
		log.Println("get letter reply homepage failed,err:", err)
		return
	}

	tool.OKWithData(ctx, letters)
}

func GetReplyLetter(ctx *gin.Context) {
	admId, res := ctx.Get("gmt_id")
	if !res {
		tool.OperateForbid(ctx)
		return
	}

	replyId, res := ctx.GetQuery("id")
	if !res {
		tool.ParamFailed(ctx)
		return
	}

	letter, err := service.GetReplyInfo(replyId, admId.(string))
	if err != nil {
		if errors.Is(err, my_consts.ErrOfForbidden) {
			tool.OperateForbid(ctx)
			return
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tool.OKWithData(ctx, "无此回复信息")
			return
		}
		tool.InternetError(ctx)
		log.Println("get letter info failed,err:", err)
		return
	}

	tool.OKWithData(ctx, letter)
}

func ChangeLetterStatus(ctx *gin.Context) {
	admId, res := ctx.Get("gmt_id")
	if !res {
		tool.OperateForbid(ctx)
		return
	}

	var req = struct {
		ReplyId int    `json:"id,omitempty"`
		IsPass  bool   `json:"is_pass"`
		Content string `json:"content,omitempty"`
	}{}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		tool.ParamFailed(ctx)
		return
	}

	err = service.ChangeLetterReplyStatus(req.ReplyId, admId.(string), req.Content, req.IsPass)
	if err != nil {
		if errors.Is(err, my_consts.ErrOfForbidden) {
			tool.OperateForbid(ctx)
			return
		}
		log.Println("can not change status,err:", err)
		tool.InternetError(ctx)
		return
	}

	tool.OK(ctx)
}

func GetReplySum(ctx *gin.Context) {
	tool.OKWithData(ctx, service.GetReplySum())
}
