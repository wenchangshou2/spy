package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wenchangshou2/Spy/core"
	"qiniupkg.com/x/errors.v7"
)

func getHTMLContext(url string) (string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", core.GetRandomUserAgent())
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		//fmt.Println("Get请求%s返回错误:%s",url,err)
		return "", errors.New("请求" + url + "错误")
	}
	if res.StatusCode != 200 {
		return "", errors.New("请求网页失败")
	}
	body := res.Body
	defer body.Close()
	bodyByte, _ := ioutil.ReadAll(body)
	resStr := string(bodyByte)
	return resStr, nil
}
func getNextURL(url string) {
	resStr, err := getHTMLContext(url)
	if err != nil { //获取网页失败
		return
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(resStr))
	if err != nil {
		fmt.Println("获取dom元素失败")
		return
	}
	dom.Find("tbody>tr").Each(func(i int, selection *goquery.Selection) {
		if i < 3 {
			return
		}
		version := selection.Find("td:nth-child(2)").Text()
		updateTime := selection.Find("td:nth-child(3)").Text()
		m, _ := regexp.MatchString("[\\d\\.]+/", version)
		if !m {
			return
		}
		spy(url + version)
		fmt.Println(updateTime, version)
	})
}

// 获取下载的链接
func getDownURL(url string) {
	resStr, err := getHTMLContext(url)
	if err != nil { //获取网页失败
		return
	}
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(resStr))
	if err != nil {
		fmt.Println("获取dom元素失败")
		return
	}
	dom.Find("tbody>tr").Each(func(i int, selection *goquery.Selection) {
		if i < 3 {
			return
		}
		downloadText := selection.Find("td:nth-child(2)").Text()
		updateTime := selection.Find("td:nth-child(3)").Text()
		fmt.Println(updateTime)
		if m,_:=regexp.MatchString(downloadText,"windows");!m{

		}
	})
}
func spy(url string) {
	fmt.Println(url)
	if len(strings.Split(url, "/")) < 8 {
		getNextURL(url)
	} else {
		getDownURL(url)
		fmt.Println("download:" + url)
	}
}
