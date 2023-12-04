/**
 * @Author: Hardews
 * @Date: 2023/4/4 23:21
 * @Description:
**/

package model

import (
	"gorm.io/gorm"
	"time"
)

type BaiduToken struct {
	gorm.Model
	ExpiresAt time.Time
	Token     string
}

// AccTokenJSON 从api获取到的响应json，只取过期时间和acc token
type AccTokenJSON struct {
	ExpiresIn   int    `json:"expires_in"`   // 过期时间
	AccessToken string `json:"access_token"` // acc token
}

// Tendency 情感倾向分析响应json
type Tendency struct {
	Text  string `json:"text"`
	Items []struct {
		Confidence   float64 `json:"confidence"`
		NegativeProb float64 `json:"negative_prob"`
		PositiveProb float64 `json:"positive_prob"`
		Sentiment    int     `json:"sentiment"`
	} `json:"items"`
	LogId int64 `json:"log_id"`
}

// ContentAudit 内容审核响应json
type ContentAudit struct {
	Conclusion string `json:"conclusion"`
	Data       []struct {
		Msg string `json:"msg,omitempty"`
	} `json:"data,omitempty"`
	ConclusionType int `json:"conclusionType"`
}
