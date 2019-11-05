package tool_crontab

import (
	"go_core/tool/tool_date"
	"time"

	"github.com/hunterhug/marmot/miner"
)

type Crontab struct {
	Year   int `json:"year"`
	Month  int `json:"month"`
	Day    int `json:"day"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}

const (
	SecondTime = 1
	MinuteTime = 60 * SecondTime
	HourTime   = 60 * MinuteTime
	DayTime    = 24 * HourTime
)

var (
	log = miner.Log()
)

//golang 定时器，启动的后按照指定时间执行
/**
目前只支持天以后的时间，
*/
func StartTimerBar(f func(), cron Crontab, isSleep bool, firstTask bool) {
	for {
		sleepTime := 0
		if firstTask {
			f()
		}
		now := time.Now()
		if isSleep {
			//睡眠一定时间之后执行任务
			sleepTime = cron.Day*DayTime + cron.Hour*HourTime + cron.Minute*MinuteTime + cron.Second*SecondTime
			next := now.Add(time.Second * time.Duration(sleepTime))
			log.Infof("sleep模式-->下次执行时间是%s", next.Format(tool_date.TimeLayout))
			t := time.NewTimer(next.Sub(now))
			<-t.C
		} else {
			//每天的几点几分几秒执行任务
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), cron.Hour, cron.Minute, cron.Second, 0, next.Location())
			log.Infof("定时模式-->下次执行时间是%s", next.Format(tool_date.TimeLayout))
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
		//执行方法
		if firstTask == false {
			f()
		}
	}
}

//明天的几点几分执行任务
func TomorrowToDo(cron Crontab)  {
	now := time.Now()
	next := now.Add(time.Hour * 24)
	next = time.Date(next.Year(), next.Month(), next.Day(), cron.Hour, cron.Minute, cron.Second, 0, next.Location())
	log.Infof("定时模式-->下次执行时间是%s", next.Format(tool_date.TimeLayout))
	t := time.NewTimer(next.Sub(now))
	<-t.C
}

func ToDayToDo(cron Crontab)  {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), cron.Hour, cron.Minute, cron.Second, 0, now.Location())
	log.Infof("定时模式-->下次执行时间是%s", next.Format(tool_date.TimeLayout))
	t := time.NewTimer(next.Sub(now))
	<-t.C
}
