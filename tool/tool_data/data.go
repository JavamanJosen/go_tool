package tool_data

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"github.com/axgle/mahonia"
	"github.com/hunterhug/marmot/miner"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	log = miner.Log()
)

func StrToInt(str string) int {
	if str == "" {
		return 0
	}
	strInt, _ := strconv.Atoi(str)
	return strInt
}

func StrToInt64(str string) int64 {
	sr, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Error(err.Error())
	}
	return sr
}

func StringToFloat(str string) (float64, error) {
	v1, err := strconv.ParseFloat(str, 64)
	return v1, err
}

func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

/**
gzip加密
*/
func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		writer.Close()
		return out, err
	}
	err = writer.Close()
	if err != nil {
		return out, err
	}

	return buffer.Bytes(), nil
}

/**
gzip解密
*/
func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

//设置编码
func ChangeEncode(html, codeType string) string {
	//	"GB18030"
	dec := mahonia.NewDecoder(codeType)
	return dec.ConvertString(string(html))
}

func UnicodeToHz(code string) string {
	bs, err := hex.DecodeString(strings.Replace(code, `\u`, ``, -1))
	if err != nil {
		return ""
	}
	context := ""
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		context += string(r)
	}
	return context
}

/**
返回字符个数
*/
func Length(str string) int {
	var r = []rune(str)
	return len(r)
}

/**
中文字符在整个文中的位置
*/
func IndexOf(source, tag string) int {
	var r = []rune(source)
	var t = []rune(tag)
	reIndex := -1

	for index, value := range r {
		count := 0
		if value == t[0] {
			count++
			for index2, item := range t {
				if index2 == 0 {
					continue
				}
				if index+index2 < len(r) && loopSearch(r[index+index2], item) {
					count++
					reIndex = index
				}
			}
		}
		if count == len(t) {
			return reIndex
		}
	}
	return reIndex
}

func loopSearch(s1, s2 int32) bool {
	if s1 == s2 {
		return true
	}
	return false
}

func Substring(source string, start, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}

func IsRegex(s1, re string) bool {
	bo, _ := regexp.MatchString(re, s1)
	if bo {
		return true
	}
	return false
}

/**
字符串md5加密
*/
func Md5(params ...string) string {
	str := ""
	for _, item := range params {
		str += item
	}
	m := md5.New()
	m.Write([]byte(str))
	result := m.Sum(nil)
	return hex.EncodeToString(result)
}

func GetHash(str string) string {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(str))
	return string(Sha1Inst.Sum([]byte("")))
}

func ContainsItem(arr []string, str string) (int,bool) {
	for index, item := range arr {
		if item  == str{
			return index, true
		}
	}
	return -1, false
}

func SplitItem(arr []string, index int) []string {
	return append(arr[:index], arr[index+1:]...)
}
func RemoveItem(arr []string, str string) []string {
	index, b := ContainsItem(arr, str)
	if b && index > -1{
		return append(arr[:index], arr[index+1:]...)
	}
	return arr
}

/**
保存内容到文件，追加的形式
*/
func Savefile(path, str_content string) {
	fd, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	//fd_time := time.Now().Format("2006-01-02 15:04:05")
	fd_content := strings.Join([]string{str_content, "\n"}, "")
	buf := []byte(fd_content)
	fd.Write(buf)
	fd.Close()
}

/**
 * 读文件内容
 */
func ReadFile(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreatePath(path string) error {
	return os.Mkdir(path, os.ModePerm)
}

func DeleteFile(path string) error {
	return os.Remove(path)
}

func AddAll(sum, sub []string) []string {
	if len(sub) == 0{
		return sum
	}
	for _, su := range sub{
		sum = append(sum, su)
	}
	return sum
}

/**
	去重
 */
func RemoveDubbo(urls []string) []string {
	reUrl := []string{}
	if len(urls) == 0{
		return reUrl
	}

	set := make(map[string]interface{})
	for _, item := range urls{
		set[item] = nil
	}

	for v := range set{
		reUrl = append(reUrl, v)
	}
	return reUrl
}

const pattern = "\\d+" //反斜杠要转义
//判断字符串是否全为数字
func IsNum(str string) bool {
	result,_ := regexp.MatchString(pattern,str)
	return result
}
