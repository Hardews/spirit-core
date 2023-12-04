/**
 * @Author: Hardews
 * @Date: 2023/4/9 20:30
 * @Description:
**/

package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"spirit-core/my_consts"
	"spirit-core/service"
	"spirit-core/tool"
)

func CheckCode(ctx *gin.Context) {
	addr, exist := ctx.GetPostForm("address")
	if !exist {
		tool.ParamFailed(ctx)
		return
	}

	code, exist := ctx.GetPostForm("code")
	if !exist {
		tool.ParamFailed(ctx)
		return
	}

	err := service.CheckCode(addr, code)
	if err != nil {
		if errors.Is(err, my_consts.ErrOfWrongCode) {
			tool.OKWithData(ctx, err.Error())
			return
		}
		log.Println("check code error,err:", err)
		tool.InternetError(ctx)
		return
	}

	tool.OK(ctx)
}

func SendCode(ctx *gin.Context) {
	addr, exist := ctx.GetPostForm("address")
	if !exist {
		tool.ParamFailed(ctx)
		return
	}

	err := service.SendCode(addr)
	if err != nil {
		tool.InternetError(ctx)
		return
	}

	tool.OK(ctx)
}
