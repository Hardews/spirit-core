/**
 * @Author: Hardews
 * @Date: 2023/3/9 23:07
 * @Description:
**/

package api

import (
	"errors"
	"spirit-core/service"
	"spirit-core/tool"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"log"
)

func UpdateLetterStatus(ctx *gin.Context) {
	var req struct {
		Id     int `json:"id,omitempty"`
		Status int `json:"status,omitempty"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		tool.ParamFailed(ctx)
		return
	}

	err = service.UpdateLetterStatus(req.Id, req.Status)
	if err != nil {
		tool.OKWithData(ctx, "有人正在操作，请稍后")
		log.Println("update letter status failed,err:", err)
		return
	}

	tool.OK(ctx)
}

func GetLetterStatus(ctx *gin.Context) {
	var err error
	id, res := ctx.GetQuery("id")
	if !res {
		tool.ParamFailed(ctx)
		return
	}

	var resp struct {
		Status int `json:"status"`
	}

	resp.Status, err = service.GetLetterStatus(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tool.OKWithData(ctx, "无此信件")
			return
		}
		tool.InternetError(ctx)
		log.Println("get letter info failed,err:", err)
		return
	}

	tool.OKWithData(ctx, resp)
}

func AddReplyLetter(ctx *gin.Context) {
	admId, res := ctx.Get("gmt_id")
	if !res {
		tool.OperateForbid(ctx)
		return
	}

	var req = struct {
		Id      int    `json:"id,omitempty"`
		Content string `json:"content,omitempty"`
	}{}

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		tool.ParamFailed(ctx)
		return
	}

	err = service.ReplyLetter(req.Id, admId.(string), req.Content)
	if err != nil {
		tool.InternetError(ctx)
		log.Printf("reply letter failed,letterId:%d, admId:%s, err:%s\n", req.Id, admId, err)
		return
	}

	tool.OK(ctx)
}

func GetLetter(ctx *gin.Context) {
	letterId, res := ctx.GetQuery("id")
	if !res {
		tool.ParamFailed(ctx)
		return
	}

	letter, err := service.GetLetterInfo(letterId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tool.OKWithData(ctx, "无此信件")
			return
		}
		tool.InternetError(ctx)
		log.Println("get letter info failed,err:", err)
		return
	}

	tool.OKWithData(ctx, letter)
}

func GetLetterHomepage(ctx *gin.Context) {
	letters, err := service.GetLetterHomepage()
	if err != nil {
		tool.InternetError(ctx)
		log.Println("get letter homepage failed,err:", err)
		return
	}

	tool.OKWithData(ctx, letters)
}
