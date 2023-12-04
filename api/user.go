/**
 * @Author: Hardews
 * @Date: 2023/4/5 21:57
 * @Description:
**/

package api

import (
	"errors"
	"spirit-core/model"
	"spirit-core/my_consts"
	"spirit-core/service"
	"spirit-core/tool"

	"github.com/gin-gonic/gin"

	"log"
)

func Register(ctx *gin.Context) {
	gmtId, res := ctx.Get("gmt_id")
	if !res {
		tool.OperateForbid(ctx)
		return
	}
	var user = model.User{}
	var exist bool
	user.Username, exist = ctx.GetPostForm("username")
	if !exist {
		tool.ParamFailed(ctx)
		return
	}
	user.Password, exist = ctx.GetPostForm("password")
	if !exist {
		tool.ParamFailed(ctx)
		return
	}
	user.Nickname, exist = ctx.GetPostForm("nickname")
	if !exist {
		tool.ParamFailed(ctx)
		return
	}
	user.Address, exist = ctx.GetPostForm("address")
	if !exist {
		tool.ParamFailed(ctx)
		return
	}

	tool.OKWithData(ctx, service.AddUser(gmtId.(string), user))

}

func Token(ctx *gin.Context) {
	var req struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil || req.Username == "" || req.Password == "" {
		log.Println("get param failed,err:", err)
		tool.ParamFailed(ctx)
		return
	}

	var user = model.User{
		Username: req.Username,
		Password: req.Password,
	}

	// 登陆
	token, rt, err := service.GetToken(user)
	if err != nil {
		if errors.Is(err, my_consts.ErrOfNoThisUser) || errors.Is(err, my_consts.ErrOfWrongPassword) {
			tool.OKWithData(ctx, err.Error())
			return
		} else {
			log.Println("get token failed,err:", err)
			tool.InternetError(ctx)
			return
		}
	}

	resp := struct {
		Token        string `json:"token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}{
		token,
		rt,
	}

	tool.OKWithData(ctx, resp)
}

func TokenRefresh(ctx *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}

	err := ctx.ShouldBindJSON(&req)
	if err != nil || req.RefreshToken == "" {
		log.Println("get param failed,err:", err)
		tool.ParamFailed(ctx)
		return
	}

	if len(req.RefreshToken) != 64 {
		tool.NeedLogin(ctx)
		return
	}

	token, rt, err := service.RefreshToken(req.RefreshToken)
	if err != nil {
		log.Println("get token failed,err:", err)
		tool.InternetError(ctx)
		return
	}

	resp := struct {
		Token        string `json:"token,omitempty"`
		RefreshToken string `json:"refresh_token,omitempty"`
	}{
		token,
		rt,
	}

	tool.OKWithData(ctx, resp)
}

func ChangePasswordByToken(ctx *gin.Context) {
	var user model.User
	iUsername, _ := ctx.Get("username")
	user.Username = iUsername.(string)

	ChangePassword(ctx, user)
}

func ChangePassword(ctx *gin.Context, user model.User) {
	var res bool
	gmtId, res := ctx.Get("gmt_id")
	if !res {
		tool.OperateForbid(ctx)
		return
	}

	// 检验入参
	user.Password, res = ctx.GetPostForm("new")
	if !res {
		tool.ErrorWithData(ctx, "新密码为空")
		return
	}

	oldPassword, res := ctx.GetPostForm("old")
	if !res {
		tool.ErrorWithData(ctx, "旧密码为空")
		return
	}

	res, err := service.ChangePassword(user, oldPassword, gmtId.(string))
	if !res {
		tool.InternetError(ctx)
		log.Println(":change password failed,err:" + err.Error())
		return
	}
	if err != nil {
		tool.ErrorWithData(ctx, err.Error())
		return
	}

	tool.OK(ctx)
}

func GetIdentity(ctx *gin.Context) {
	gmtId, exist := ctx.Get("gmt_id")
	if !exist {
		tool.OperateForbid(ctx)
		return
	}

	tool.OKWithData(ctx, service.GetIdentity(gmtId.(string)))
}
