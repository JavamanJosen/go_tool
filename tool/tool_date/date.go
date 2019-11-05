package tool_date

import (
	"go_core/tool/tool_data"
	"fmt"
	"time"
)

const (
	//转化所需模板
	TimeLayout = "2006-01-02 15:04:05"
)

/**
时间工具类
*/

//获取几天前或者几天后的时间  -1或1
func GetDayTime(n int) string {
	currTime := time.Now()
	return currTime.AddDate(0, 0, n).Format("2006-01-02 15:04:05")
}

//获取当前时间点的格式化时间
func GetNowTimeFormat() string {
	currTime := time.Now()
	return currTime.Format("2006-01-02 15:04:05")
}

/**
入参是以当前为基准 昨天到明天  s1=-1， e1 = 1，t指定返回类型，
t=format返回格式化时间yyyy-MM-dd HH:mm:ss,其他范围时间戳
isNatural 是否是自然天， 为true时 结果是 2019-04-01 00:00:00
*/
func GetStartEndTime(s1, e1 int, t string, isNatural bool) (string, string) {
	st := GetDayTime(s1)
	en := GetDayTime(e1)
	if isNatural {
		st = st[:10] + " 00:00:00"
		en = en[:10] + " 23:59:59"
	}

	if t == "format" {
		return st, en
	}
	return fmt.Sprintln(FormatTimeToInt(st)), fmt.Sprintln(FormatTimeToInt(en))

}

func GetStartEndTimeInt(s1, e1 int, isNatural bool) (int64, int64) {
	st := GetDayTime(s1)
	en := GetDayTime(e1)
	if isNatural {
		st = st[:10] + " 00:00:00"
		en = en[:10] + " 23:59:59"
	}
	return FormatTimeToInt(st), FormatTimeToInt(en)

}

func FormatTimeToInt(datetime string) int64 {
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(TimeLayout, datetime, loc)
	return tmp.Unix() //转化为时间戳 类型是int64
}

/**
时间戳转时间格式
*/
func FormatToString(t string) string {
	sr := tool_data.StrToInt64(t)
	//时间戳转日期
	dataTimeStr := time.Unix(sr, 0).Format(TimeLayout) //设置时间戳 使用模板格式化为日期字符串
	return dataTimeStr
}

func GetTodayFormat(n int) string {
	return GetDayTime(n)[:10]
}

func FormatInt64ToStr(t int64) string {
	dataTimeStr := time.Unix(t, 0).Format(TimeLayout) //设置时间戳 使用模板格式化为日期字符串
	return dataTimeStr
}

func GetNowDate() string {
	return time.Now().Format("20060102")
}

func GetNowTime() string {
	return time.Now().Format(TimeLayout)
}

//根据起止时间获取之间的日期
//GetTimeListByStartEnd("2019-08-01", "2019-08-05", true, true)
func GetTimeListByStartEnd(start, end string, paramIsFormat, resultIsFormat bool) []string {
	reList := []string{}
	var startInt int64
	var endInt int64
	if paramIsFormat{
		if len(start) == 10{
			start = start + " 00:00:00"
		}
		if len(end) == 10{
			end = end + " 00:00:00"
		}
		startInt = FormatTimeToInt(start)
		endInt = FormatTimeToInt(end)
	}else{
		startInt = tool_data.StrToInt64(start)
		endInt = tool_data.StrToInt64(end)
	}

	if endInt<startInt{
		return reList
	}

	for i:=endInt; i>= startInt; i -= 86400{
		var t string
		if resultIsFormat{
			t = FormatInt64ToStr(i)
		}else{
			t = tool_data.Int64ToString(i)
		}
		reList = append(reList, t)
	}
	return reList
}
