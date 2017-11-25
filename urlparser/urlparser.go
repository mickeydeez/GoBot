package urlparser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type UrlParser struct {
	url     string
	content []byte
}

func InitParser(url string) UrlParser {
	parser := UrlParser{}
	parser.url = url
	resp, err := http.Get(url)
	if checkError(err) {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if checkError(err) {
			parser.content = body
		}
	}
	return parser
}

func (u UrlParser) RootDomain() string {
	split := strings.Split(u.url, "/")
	return split[2]
}

func (u UrlParser) Title() string {
	regex, err := regexp.Compile("<title>(.*)</title>")
	status := checkError(err)
	var title string
	if status {
		matches := regex.FindSubmatch(u.content)
		if len(matches) != 2 {
			title = ""
		} else {
			title = string(matches[1])
		}
	}
	return title
}

func checkError(err error) bool {
	if err != nil {
		fmt.Println("Critical Error: ", err)
		return false
	}
	return true
}
