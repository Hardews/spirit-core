/**
 * @Author: Hardews
 * @Date: 2023/4/5 21:59
 * @Description:
**/

package tool

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"spirit-core/dao"
)

type GMTClaims struct {
	AdmID                string `json:"admID,omitempty"`
	jwt.RegisteredClaims        // v4版本新增，原jwt.StandardClaims
}

const (
	salt              = "gmt-website"
	refreshExpiration = 60 * 60 * 24 * 14 * time.Second // refresh token 过期时间，设为14天
)

var GMTSecret = []byte("gm_team")

var ErrOfRefresh = errors.New("refresh token failed")

// TokenGenerate 生成token
func TokenGenerate(admId string) (acToken string, reToken string, err error) {
	claim := GMTClaims{
		AdmID: admId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                       // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                       // 生效时间
		}}

	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法

	acToken, err = Token.SignedString(GMTSecret)

	reToken = refreshTokenGenerate(acToken, admId)

	return acToken, reToken, err
}

// TokenVerify 验证token
func TokenVerify(accToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accToken, &GMTClaims{}, secret())
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", errors.New("that's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", errors.New("token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "", errors.New("token not active yet")
			} else {
				return "", errors.New("couldn't handle this token")
			}
		}
	}

	if claims, ok := token.Claims.(*GMTClaims); ok && token.Valid {
		return claims.AdmID, nil
	}

	return "", errors.New("couldn't handle this token")
}

// RefreshToken 刷新token
func RefreshToken(rt string) (acToken string, reToken string, err error) {
	// 先从redis拿
	gmtId, err := Get(rt)
	if len(gmtId) != 64 {
		err = ErrOfRefresh
	}
	if err != nil {
		log.Println("get rt from redis failed, err:", err)

		// 从mysql拿
		gmtId = dao.GetRefreshToken(rt)
	}

	if gmtId == "" {
		err = ErrOfRefresh
		return
	}

	go Del(rt)
	go dao.DelRefreshToken(rt, gmtId)

	return TokenGenerate(gmtId)
}

// 生成方案:取payload + salt 进行sha256加密
func refreshTokenGenerate(accToken, gmtId string) string {
	token := strings.Split(accToken, ".")
	payload := token[1]

	// 生成refresh token
	h := sha256.New()
	h.Write([]byte(payload + salt))
	refreshToken := fmt.Sprintf("%x", h.Sum(nil))

	// 存储,redis应是弱依赖，异步存储
	go Set(refreshToken, gmtId, refreshExpiration)
	err := dao.AddRefreshToken(refreshToken, gmtId)
	if err != nil {
		log.Printf("set refresh token failed,gmt_id:%s, rt:%s, err:%s\n", gmtId, refreshToken, err)
		refreshToken = ""
	}

	return refreshToken
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return GMTSecret, nil
	}
}
