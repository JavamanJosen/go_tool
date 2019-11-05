package request

import (
	"bytes"
	"github.com/hunterhug/marmot/miner"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	contentType = "application/json;charset=utf-8"
)

var log = miner.Log()
type SRequest struct {
	Url      string
	Method   string
	Headers  map[string][]string
	Body     string
	BodyByte []byte
	Form     map[string][]string
	Cookie   []string
	Referer  string
	Timeout  time.Duration
	ReTry    int
}

func Get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)
	return string(re), err
}

func GetHeadersData(requestUrl string, headers, data map[string][]string) (string, error) {
	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(url.Values.Encode(data)))
	if err != nil {
		println(err.Error())
	}
	req.Header = headers
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)
	return string(re), err
}

func GetHeaders(url string, headers map[string][]string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err.Error())
	}
	req.Header = headers
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)
	return string(re), err
}

func GetHeadersRespCookie(url string, headers map[string][]string) (string, []string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		println(err.Error())
	}
	req.Header = headers
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println(err.Error())
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)

	coo := []string{}
	for _, c := range resp.Cookies() {
		coo = append(coo, c.Name+"="+c.Value)
	}
	return string(re), coo
}

func PostJson(url string, b []byte) (string, error) {
	body := bytes.NewBuffer(b)
	resp, err := http.Post(url, contentType, body)
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)
	return string(re), err
}

func PostForm(url string, data map[string][]string) (string, error) {
	resp, err := http.PostForm(url, data)
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}
	defer resp.Body.Close()
	re, err := ioutil.ReadAll(resp.Body)
	return string(re), err
}

func GetHtmlByResp(resp *http.Response) string {
	re, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf(err.Error())
		return ""
	}
	return string(re)
}

func GetCookieByResp(resp *http.Response) []string {
	cookie := []string{}
	for _, c := range resp.Cookies() {
		cookie = append(cookie, c.Name+"="+c.Value)
	}
	return cookie
}
