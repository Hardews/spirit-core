/**
 * @Author: Hardews
 * @Date: 2023/3/9 17:51
 * @Description:
**/

package api

import (
	"spirit-core/tool"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()

	r.Use(tool.Cors())

	// 管理员登录相关
	r.POST("/token", Token)
	r.POST("/token/refresh", TokenRefresh)

	// 管理员信息相关（如密码
	user := r.Group("/user")
	{
		user.Use(tool.Verify())
		user.PUT("/password", ChangePasswordByToken) // 修改密码
	}

	// 信件相关
	letter := r.Group("/letter")
	{
		letter.Use(tool.AntiHarassment())
		code := letter.Group("/code")
		{
			code.POST("", SendCode)
			code.POST("/check", CheckCode)
		}
		letter.POST("", AddLetter)            // 写信
		letter.GET("/sum", GetLetterSum)      // 获取信件的总和
		letter.GET("/reply/sum", GetReplySum) // 获取回复信件的总和

		future := letter.Group("/future")
		{
			future.POST("", AddFutureLetter)       // 给未来的自己写信
			future.GET("", GetFutureLetter)        // 获取公开信
			future.GET("/id", GetFutureLetterById) // 获取一封公开信的详情
			future.GET("/random", GetRandomFutureLetter)
			future.GET("/sum", GetFutureLetterSum)
		}
	}

	// 树洞相关接口
	hole := r.Group("/hole")
	{
		hole.GET("/new", NewTreeHole)
		hole.GET("/join", JoinTreeHole)
	}

	// 信件管理员相关
	adm := r.Group("/admin")
	{

		adm.Use(tool.Verify())
		adm.GET("/info", GetIdentity) // 获取管理员的身份信息

		letterAdm := adm.Group("/letter")
		{
			// 处理信件相关
			letterAdm.GET("", GetLetter)                  // 获取某个信件的详细信息
			letterAdm.POST("", AddReplyLetter)            // 回复信件
			letterAdm.GET("/status", GetLetterStatus)     // 获取信件状态
			letterAdm.POST("/status", UpdateLetterStatus) // 修改信件状态
			letterAdm.GET("/homepage", GetLetterHomepage) // 获取信件概览

			// 处理信件回复相关
			reply := letterAdm.Group("/reply")
			{
				// 普通管理员
				reply.GET("", GetMyReplyOverview)          // 获取我回复的信件概览
				reply.GET("/advice", GetReplyLetterAdvice) // 获取回复未通过的原因

				// 超级管理员
				reply.GET("/homepage", GetLetterReplyHomepage) // 获取待审核信件的概览
				reply.GET("/status", GetReplyLetter)           // 获取待审核的信件
				reply.PUT("/status", ChangeLetterStatus)       // 修改信件的状态（回复信件的内容是否通过审核
				reply.POST("/register", Register)              // 注册
			}
		}
	}

	r.Run(":8090")
}
