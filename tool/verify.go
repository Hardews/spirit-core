/**
 * @Author: Hardews
 * @Date: 2023/4/5 22:00
 * @Description:
**/

package tool

import "github.com/gin-gonic/gin"

func Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		gmtId, err := TokenVerify(token)
		if err != nil {
			ErrorWithData(ctx, err.Error())
			ctx.Abort()
			return
		}

		ctx.Set("gmt_id", gmtId)
		ctx.Next()
	}
}

func AntiHarassment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		if token != "gmt-website" {
			OperateForbid(ctx)
			ctx.Abort()
			return
		}
	}
}
