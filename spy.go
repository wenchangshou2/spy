package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	//"time"

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
func getDownloadFileMd5(url string,name string)(string,error){
	c,err:=getHTMLContext(url+"md5sums.txt")
	if err!=nil{
		return "",errors.New("下载失败")
	}
	scanner:=bufio.NewScanner(strings.NewReader(c))
	for scanner.Scan(){
		txt:=scanner.Text()
		if strings.Index(txt,name)>-1{
			idx:=strings.Index(txt," ")
			fmt.Println(txt,txt[0:idx])
			return txt[0:idx],nil
		}
	}
	return "",errors.New("没有找到对应的内容")

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
		//updateTime := selection.Find("td:nth-child(3)").Text()


		if i:=strings.Index(downloadText,"exe");i>-1{
			urlList:=strings.Split(url,"/")
			m,_:=getDownloadFileMd5(url,downloadText)
			d:=downloadContent{}
			d.MainVersion=urlList[len(urlList)-3]
			d.SubVersion=urlList[len(urlList)-2]
			d.Md5=m
			d.DownloadUrl=url+downloadText
			downloadUrlList=append(downloadUrlList, d)
		}
	})
}

func isIgore(url string)bool{
	list:=strings.Split(url,"/")
	version:=list[len(list)-2]
	if m,_:=regexp.MatchString("\\d+",version);m{
		for _,v:=range ignoreList{
			if strings.Compare(v,version) ==0{
				return true
			}

		}
	}
	return false
}
func spy(url string) {
	fmt.Println(url)
	if isIgore(url){
		return
	}
	//return
	if len(strings.Split(url, "/")) < 8 {
		getNextURL(url)
	} else {
		getDownURL(url)
		fmt.Println("download:" + url)
	}
}
