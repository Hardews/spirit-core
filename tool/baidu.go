/**
 * @Author: Hardews
 * @Date: 2023/4/4 23:17
 * @Description:baidu api相关接口
**/

package tool

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"spirit-core/dao"
	"spirit-core/model"
)

var mu sync.Mutex

var (
	clientId     = os.Getenv("spirit_core_baidu_client_id")
	clientSecret = os.Getenv("spirit_core_baidu_client_secret")
)

// getBaiduAccToken 获取百度api所需的acc token
func getBaiduAccToken() (string, error) {
	// 先从数据库拿,拿不到再申请
	accToken, err := dao.GetAccToken()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			goto get
		}
		return "", err
	}

	// 获取到了没过期就返回
	if accToken.Token != "" && accToken.ExpiresAt.After(time.Now()) {
		return accToken.Token, nil
	}

get:
	// 没获取到就发送请求获取,加锁防并发
	if mu.TryLock() {
		defer mu.Unlock()
		url := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&" +
			"client_id=" + clientId +
			"&client_secret=" + clientSecret

		client := &http.Client{}
		var req *http.Request
		var res *http.Response
		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			return "", err
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		res, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}

		// 获取到新的token之后,添加后返回
		var accTokenJson model.AccTokenJSON
		err = json.Unmarshal(body, &accTokenJson)
		if err != nil {
			return "", err
		}

		// 提前一天判过期
		err = dao.AddAccToken(model.BaiduToken{
			ExpiresAt: time.Now().AddDate(0, 0, accTokenJson.ExpiresIn/(60*24)-1),
			Token:     accTokenJson.AccessToken,
		})
		if err != nil {
			// 不行就重新搞
			log.Println("get baidu token failed, err:", err)
			return getBaiduAccToken()
		}

		return accTokenJson.AccessToken, nil
	} else {
		// 睡眠后尝试再次获取
		time.Sleep(time.Second)
		return getBaiduAccToken()
	}
}

// AffectiveTendencyAnalysis 情感倾向分析
func AffectiveTendencyAnalysis(text string) (model.Tendency, error) {
	accToken, err := getBaiduAccToken()
	if err != nil {
		log.Println(err)
		return model.Tendency{}, err
	}

	url := "https://aip.baidubce.com/rpc/2.0/nlp/v1/sentiment_classify?access_token=" + accToken + " &charset=UTF-8"

	body, err := sendPostByJsonForm(url, text)
	log.Println(err)

	var ten model.Tendency
	err = json.Unmarshal(body, &ten)

	return ten, err
}

// ContentAudit 内容审核
func ContentAudit(text string) (model.ContentAudit, error) {
	accToken, err := getBaiduAccToken()
	if err != nil {
		return model.ContentAudit{}, err
	}

	url := "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined?access_token=" + accToken + "&charset=UTF-8"

	body, err := sendPostByWWW(url, map[string]string{"text": text})

	log.Println(string(body))
	var aud model.ContentAudit
	err = json.Unmarshal(body, &aud)

	return aud, err
}
