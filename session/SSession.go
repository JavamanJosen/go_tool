package session

import (
	"go_core/request"
	"go_core/response"
	"go_core/tool/tool_parser"
	"bytes"
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Timeout     = 10
	contentType = "application/json;charset=utf-8"
	proxyServer = ""
)

var (
	log = miner.Log()
)

type Session struct {
	Request  request.SRequest
	Response response.SResponse
	Cookie   []string
	UseProxy bool
	IsVerify bool
	Client   http.Client
	AppName  string
}

type AbuyunProxy struct {
	AppID     string
	AppSecret string
}

func (p AbuyunProxy) ProxyClient() *http.Client {
	proxyUrl, _ := url.Parse("http://" + p.AppID + ":" + p.AppSecret + "@" + proxyServer)
	return &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
}

func GetSession(useProxy bool) *Session {
	session := &Session{}
	//session 的 request
	headers := map[string][]string{}
	headers["User-Agent"] = []string{miner.RandomUa()}
	session.Request.Headers = headers
	session.Request.Timeout = Timeout * time.Second
	session.UseProxy = useProxy
	//session的response
	session.Response.StatusCode = 200
	session.Client = http.Client{}

	return session
}

func (session *Session) Send() (*Session, error) {
	var err error
	req := &http.Request{} //用来接收http.request
	if (session.Request.Method == "POST" || session.Request.Method == "post") && session.Request.Body != "" {
		req, err = http.NewRequest(session.Request.Method, session.Request.Url, bytes.NewBuffer([]byte(session.Request.Body)))
		req.Header.Set("Content-Type", contentType)
	} else if (session.Request.Method == "POST" || session.Request.Method == "post") && session.Request.BodyByte != nil {
		req, err = http.NewRequest(session.Request.Method, session.Request.Url, bytes.NewReader(session.Request.BodyByte))
		req.Header.Set("Content-Type", contentType)
	} else if session.Request.Method == "POST" || session.Request.Method == "post" {
		req, err = http.NewRequest(session.Request.Method, session.Request.Url, strings.NewReader(url.Values.Encode(session.Request.Form)))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if "GET" == session.Request.Method || "get" == session.Request.Method {
		req, err = http.NewRequest(session.Request.Method, session.Request.Url, nil)
	} else {
		req, err = http.NewRequest(session.Request.Method, session.Request.Url, nil)
	}

	if err != nil {

		log.Error(err.Error())
		return session, err
	}

	//拼接referer
	if session.Request.Referer != "" {
		req.Header["Referer"] = []string{session.Request.Referer}
	}

	//拼接header
	headers := session.Request.Headers
	for key, _ := range headers {
		if strings.Contains(key, "Cookie") {
			continue
		}
		req.Header[key] = headers[key]
	}

	//拼接cookie
	if len(session.Cookie) > 0 {
		for _, cookie := range session.Cookie {
			cookieArr := strings.Split(cookie, "=")
			cookie1 := &http.Cookie{Name: cookieArr[0], Value: cookieArr[1], HttpOnly: true}
			req.AddCookie(cookie1)
		}
	}

	client := session.Client

	if session.UseProxy {//获取代理
		proxyUser, proxyPass := GetProxy(session.AppName)
		client = *AbuyunProxy{proxyUser, proxyPass}.ProxyClient()
	}

	client.Timeout = session.Request.Timeout

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err.Error())
		return session, err
	}

	defer resp.Body.Close()

	if session.IsVerify {
		byteArr, err := session.GetByteByRespBody(resp.Body)
		if err != nil {
			log.Error(err.Error())
		}
		session.Response.Byte = byteArr
		session.IsVerify = false
	}
	session.Response.StatusCode = resp.StatusCode
	session.Response.Status = resp.Status
	session.Response.Html = GetHtml(resp)
	session.Cookie = getCookies(session.Cookie, resp.Cookies())
	return session, nil
}

func getHeader(lastHeader map[string][]string, newHeader map[string][]string) map[string][]string {
	for key, _ := range newHeader {
		lastHeader[key] = newHeader[key]
	}
	return lastHeader
}

/**
整合新返回的cookie和请求时候的cookie
*/
func getCookies(lastCookie []string, reCookies []*http.Cookie) []string {
	newCookies := []string{}
	cooMap := map[string]string{}
	//先把上一个链接的cookie放入map中
	for _, coo := range lastCookie {
		cooArr := strings.Split(coo, "=")
		if len(cooArr) == 2 {
			cooMap[cooArr[0]] = cooArr[1]
		}
	}
	//再把新的链接的cookie放入map
	for _, coo := range reCookies {
		cooMap[coo.Name] = coo.Value
	}
	//取出cookie
	for coo, _ := range cooMap {
		newCookies = append(newCookies, coo+"="+cooMap[coo])
	}
	return newCookies
}

/**
整合返回的header和请求时候的header
*/
func GetHtml(response *http.Response) string {
	html, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err.Error())
	}
	return string(html)
}

/**
io.reader 转 []byte  通常用于验证吗
*/
func (session *Session) GetByteByRespBody(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}



const (
	timeOut = 5 * time.Second
	//获取代理的服务
	getProxyUrl     = ""
	reTry   = 3
)

func GetProxy(appName string) (string, string) {
	sess := GetSession(false)
	sess.Request.Url = fmt.Sprintf(getProxyUrl, appName)
	sess.Request.Timeout = timeOut
	for i := 0; i < reTry; i++ {
		reSess, err := sess.Send()
		if err != nil {
			log.Error(err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		proxyUser := tool_parser.RegexStr(reSess.Response.Html, `"proxyUser":"(.*?)"`)
		proxyPass := tool_parser.RegexStr(reSess.Response.Html, `"proxyPass":"(.*?)"`)
		if proxyUser != "" && proxyPass != ""{
			return proxyUser, proxyPass
		}
	}
	return "", ""
}