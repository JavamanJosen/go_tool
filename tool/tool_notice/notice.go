package tool_notice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	notifyURL = "http://10.2.17.32:8888/%s"
	log       = miner.Log()
)

type NotifyEntity struct {
}

type MailText struct {
	Receiver string `json:"receiver"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	App      string `json:"app"`
}

func (m *MailText) Dumps() string {
	js, err := json.Marshal(m)
	if err != nil {
		log.Warning(err)
		return ""
	}
	return string(js)
}

type ImText struct {
	Owner string `json:"owner"`
	Group string `json:"group"`
	Msg   string `json:"msg"`
	App   string `json:"app"`
}

func (i *ImText) Dumps() string {
	js, err := json.Marshal(i)
	if err != nil {
		log.Warning(err)
		return ""
	}
	return string(js)
}

func (n *NotifyEntity) SendMail(receivers []string, Title, Content, App string) {
	text := MailText{
		strings.Join(receivers, ","),
		Title,
		Content,
		App,
	}
	msg := text.Dumps()
	url := fmt.Sprintf(notifyURL, "mail")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Infof("response Body: %s", string(body))
}

func (n *NotifyEntity) SendIm(owner, group, msg, app string) {
	text := ImText{
		owner,
		group,
		msg,
		app,
	}
	url := fmt.Sprintf(notifyURL, "im")
	msg = text.Dumps()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Infof("response Body: %s", string(body))

}
