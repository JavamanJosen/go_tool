package tool_parser

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
)

func Goquery(html, query, mother string) []string {
	var result []string

	node, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	node.Find(query).Each(func(_ int, elem *goquery.Selection) {

		if mother == "href" {

			if h, ok := elem.Attr("href"); ok {
				result = append(result, h)
			}

		} else if mother == "text" {

			result = append(result, elem.Text())

		} else if mother == "src" {

			if src, ok := elem.Attr("src"); ok {
				result = append(result, src)
			}

		} else {

			if h, ok := elem.Html(); ok == nil {
				result = append(result, h)
			}

		}
	})

	return result
}

func RegexParses(str, pattern string) ([]string, error) {
	var result []string

	regex := regexp.MustCompile(pattern)
	find := regex.FindAllStringSubmatch(str, -1)

	if len(find) == 0 {
		return result, errors.New("not found")
	}

	for x, _ := range find {
		result = append(result, find[x][1])
	}

	return result, nil
}

// RegexParse : 通过正则表达式提取 html中的指定 regex 元素
func Regex(html, rex string) []string {
	regex := regexp.MustCompile(rex)
	find := regex.FindAllStringSubmatch(html, -1)

	if len(find) == 0 || len(find[0]) <= 1 {
		return []string{}
	}

	return []string{find[0][1]}
}

// RegexParse : 通过正则表达式提取 html中的指定 regex 元素
func RegexStr(html, rex string) string {
	regex := regexp.MustCompile(rex)
	find := regex.FindAllStringSubmatch(html, -1)

	if len(find) == 0 || len(find[0]) <= 1 {
		return ""
	}
	if len([]string{find[0][1]}) == 0 {
		return ""
	}

	return []string{find[0][1]}[0]
}

func IsRegex(str, regex string) bool {

	flag, _ := regexp.MatchString(regex, str)

	return flag
}
