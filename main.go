package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)
type config struct{
	File struct{
		Path string `yaml:"path"`
	}
	List[] string `yaml:"ignore"`
}
type downloadContent struct{
	MainVersion string
	SubVersion string
	DownloadUrl string
	Md5 string
}
var file_save_path string
var ignoreList  []string
var downloadUrlList []downloadContent
func init() {
	c:=new(config)
	yamlFile, err := ioutil.ReadFile("config.yaml")

	//log.Println("yamlFile:", string(yamlFile))
	if err != nil {
		log.Printf("yamlFile.Get err #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err!=nil{
		log.Fatal("解析失败")
	}
	file_save_path=c.File.Path
	ignoreList=c.List
	log.Println(ignoreList)

}
func createDirectory()  {//创建目录
	if _,err:=os.Stat(file_save_path);os.IsNotExist(err){
		os.MkdirAll(file_save_path,os.ModePerm)
	}
}
func main() {
	createDirectory()
	spy("https://download.qt.io/archive/qt/")
	fmt.Println(len(downloadUrlList))
	for k,v:=range downloadUrlList{
		fmt.Println(k,v)
	}
}
