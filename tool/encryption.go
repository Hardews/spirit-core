/**
 * @Author: Hardews
 * @Date: 2023/4/5 22:01
 * @Description:
**/

package tool

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// Encryption 加密方案
func Encryption(input string) (pwdHash string, err error) {
	pwd := []byte(input)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return
	}

	pwdHash = string(hash)
	return
}

// CheckPassword 验证密码
func CheckPassword(storagePwd string, inputPwd string) bool {
	byteHash := []byte(storagePwd)
	bytePwd := []byte(inputPwd)
	return bcrypt.CompareHashAndPassword(byteHash, bytePwd) == nil
}

// GenerateGMTId gmt_id 生成方案
func GenerateGMTId(address, username string) string {
	m := sha256.New()
	m.Write(GMTSecret)
	m.Write([]byte(username))
	m.Write([]byte(address))
	return fmt.Sprintf("%x", m.Sum(nil))
}
