/**
 * @Author: Hardews
 * @Date: 2023/4/15 14:34
 * @Description:
**/

package tool

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"spirit-core/dao"
	"spirit-core/model"
	"strings"

	"log"
	"time"
)

var retry int
var c = cron.New()

func Cron() {
	SelectFutureEmail()
	dao.GetFutureLetterPublicSum()
	entryId, err := c.AddFunc("0 0 * * *", SelectFutureEmail)
	if err != nil {
		log.Printf("send future email failed,entry id:%d,time:%s,err:%s", entryId, time.Now().String(), err.Error())
	}

	entryId, err = c.AddFunc("0 0 * * *", func() {
		dao.GetFutureLetterPublicSum()
	})
	if err != nil {
		log.Printf("send future email failed,entry id:%d,time:%s,err:%s", entryId, time.Now().String(), err.Error())
	}
	c.Start()
}

// SelectFutureEmail 搜索未来信件
func SelectFutureEmail() {
	fls, err := dao.GetTodayFutureEmail()
	if err != nil {
		log.Printf("get today email failed,err is:%s,today is:%s", err.Error(), time.Now().String())
		retry++
		if retry == 3 {
			panic(err)
			return
		}
		err = nil
		SelectFutureEmail()
	}

	for _, fl := range fls {
		spec := fmt.Sprintf("%d %d %d %d *", fl.SendTime.Minute(), fl.SendTime.Hour(), fl.SendTime.Day(), fl.SendTime.Month())
		entryId, err := c.AddFunc(spec, func() {
			SendFutureEmail(fl)
		})
		if err != nil {
			log.Printf("send future email failed,entry id:%d,time:%s,err:%s", entryId, time.Now().String(), err.Error())
		}
	}
}

func SendFutureEmail(fl model.Future) {
	theme := strings.ReplaceAll(FutureTheme, "{.theme}", fl.Theme)
	content := strings.ReplaceAll(FutureContent, "{.content}", fl.Text)
	content = strings.ReplaceAll(content, "{.name}", fl.Name)
	content = strings.ReplaceAll(content, "{.addTime}", fl.CreatedAt.Format("2006/01/02"))
	content = strings.ReplaceAll(content, "{.sendTime}", fl.SendTime.Format("2006/01/02"))

	err := SendEmail(theme, fl.Address, content, Future)
	if err != nil {
		log.Printf("send future email failed,id:%d,time:%s,err:%s", fl.ID, time.Now().String(), err.Error())
		return
	}

	err = dao.UpdateFutureStatus(fl.ID)
	if err != nil {
		log.Printf("update future email status failed,id:%d,time:%s,err:%s", fl.ID, time.Now().String(), err.Error())
		return
	}
}
