/**
 * @Author: Hardews
 * @Date: 2023/4/9 16:56
 * @Description:
**/

package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"spirit-core/service"
	"spirit-core/tool"
)

func JoinTreeHole(ctx *gin.Context) {
	holeId, exist := ctx.GetQuery("room")
	if !exist {
		tool.ParamFailed(ctx)
		return
	}

	username, exist := ctx.GetQuery("id")
	if !exist {
		username = service.RandomStr(6)
	}

	err := service.JoinRoom(ctx, holeId, username)
	if err != nil {
		tool.ErrorWithData(ctx, err.Error())
		return
	}
}

func NewTreeHole(ctx *gin.Context) {
	username, exist := ctx.GetQuery("id")
	if !exist {
		username = service.RandomStr(6)
	}

	roomId, err := service.NewRoom(ctx, username)
	if err != nil {
		tool.InternetError(ctx)
		log.Println("new room failed, err:", err)
		return
	}

	log.Println(roomId)
}
