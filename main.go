package main

import (
	"github.com/hunterhug/marmot/miner"
	"go_core/session"
)

var log = miner.Log()


func main()  {
	sess := session.GetSession(false)
	sess.AppName = "test"
	sess.Request.Url = "http://baidu.com"
	reSess, err := sess.Send()
	if err != nil{
		log.Error(err.Error())
		return
	}
	log.Info(reSess.Response.Html)
}
