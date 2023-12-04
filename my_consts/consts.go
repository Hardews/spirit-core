/**
 * @Author: Hardews
 * @Date: 2023/12/4 9:21
 * @Description:
**/

package my_consts

import (
	"errors"
	"time"
)

// service 状态

const (
	CodeExp     = 5 * 60 * time.Second // code 过期时间
	DefaultSalt = "gmt-team"
)

var (
	ErrOfWrongCode      = errors.New("wrong code")
	ErrOfHarassmentInfo = errors.New("this is a harassment req")
	ErrOfForbidden      = errors.New("operate forbidden")
	ErrOfNoThisUser     = errors.New("无此用户")
	ErrOfWrongPassword  = errors.New("密码错误")

	WriteWait            = 10 * time.Second
	PongWait             = 60 * time.Second
	PingPeriod           = (PongWait * 9) / 10
	MaxMessageSize int64 = 512
)
