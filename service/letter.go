/**
 * @Author: Hardews
 * @Date: 2023/3/27 16:22
 * @Description:
**/

package service

import (
	"math/rand"
	"spirit-core/dao"
	"spirit-core/model"
	"spirit-core/my_consts"
	"spirit-core/tool"
	"strconv"
	"time"
)

func getPredicted(ten model.Tendency) float64 {
	var sentiment = float64(ten.Items[0].Sentiment - 1)
	if sentiment == 1 {
		return ten.Items[0].PositiveProb - ten.Items[0].NegativeProb
	}

	var predicted float64
	pseudoPre := ten.Items[0].Confidence * sentiment
	if pseudoPre > 0 {
		predicted = pseudoPre*ten.Items[0].PositiveProb + pseudoPre*(-ten.Items[0].NegativeProb)
	} else {
		predicted = pseudoPre*(-ten.Items[0].PositiveProb) + pseudoPre*ten.Items[0].NegativeProb
	}

	return predicted
}

func AddLetter(addr, content string) error {
	if !tool.IsEmail(addr) {
		return tool.ErrOfNotARealEmail
	}

	if !dao.IsPermAddr(addr) {
		return my_consts.ErrOfHarassmentInfo
	}

	identificationRes, err := tool.ContentAudit(content)
	if err != nil {
		return err
	}

	tenRes, err := tool.AffectiveTendencyAnalysis(content)
	if err != nil {
		return err
	}

	var msg string
	if identificationRes.ConclusionType != 1 {
		msg = identificationRes.Data[0].Msg
	}

	letter := model.Letter{
		Context:    content,
		Address:    addr,
		Predicted:  getPredicted(tenRes),
		Conclusion: identificationRes.Conclusion,
		Msg:        msg,
	}

	letterInfo := model.LetterInfo{
		Text:         content,
		Confidence:   tenRes.Items[0].Confidence,
		NegativeProb: tenRes.Items[0].NegativeProb,
		PositiveProb: tenRes.Items[0].PositiveProb,
		Sentiment:    tenRes.Items[0].Sentiment,
		LogId:        tenRes.LogId,
	}

	go dao.UpdateAddress(addr)

	return dao.AddLetter(letter, letterInfo)
}

func AddFutureLetter(letter model.LetterFutureForm) (bool, error) {
	sendTime, _ := time.ParseInLocation("2006-01-02 15:04:05", letter.Time, time.Local)
	fl := model.Future{
		Name:        letter.Name,
		Text:        letter.Content,
		Theme:       letter.Theme,
		Address:     letter.Address,
		SendTime:    sendTime,
		IsPublic:    letter.IsPublic,
		IsSend:      false,
		ShouldAudit: true,
	}

	audit, _ := tool.ContentAudit(letter.Content)
	if audit.ConclusionType == 1 {
		fl.ShouldAudit = false
	}

	err := dao.AddFutureLetter(fl)

	return fl.ShouldAudit, err
}

func GetLetterSum() int64 {
	return dao.GetLetterSum()
}

func GetFutureLetter() []model.Future {
	return dao.GetFutureEmail()
}

func GetFutureLetterById(id string) model.Future {
	return dao.GetFutureEmailById(id)
}

func GetFutureLetterSum() int64 {
	return dao.GetLetterFutureSum()
}

func GetFutureRandomLetter() model.Future {
	rand.Seed(time.Now().Unix())
	ids := dao.GetFutureIds()
	id := ids[rand.Intn(int(dao.LetterSum))]
	return dao.GetFutureEmailById(strconv.Itoa(id))
}
