package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

//设置游标
var marker int64

//func Marker() int64 {
//	p, _ := time.Parse("2006-01-02T15:04:05Z", marker)
//	marker = p.Unix()
//	return marker
//}

func Uijmio(str string) int64 {
	p, _ := time.Parse("2006-01-02T15:04:05Z", str)
	io := p.Unix()
	return io
}

func UrlFunc(starttime, endtime int64) string {
	ApiUrl := "https://api.cloudflare.com/client/v4/zones/xxxxxxxxxxxxxxx/logs/received?"
	ApiUrl += "start=" + strconv.FormatInt(starttime, 10) + "&" + "end=" + strconv.FormatInt(endtime, 10)
	//Marker(endtime)
	//fmt.Printf(ApiUrl + "&count=10\n")
	return ApiUrl
}

//获取数据
func GetData(starttime, endtime int64) {
	//t := "2020-03-08T22:00:31Z"
	//p, _ := time.Parse("2006-01-02T15:04:05Z", t)
	//fmt.Println(p.Unix())
	client := &http.Client{}
	req, err := http.NewRequest("GET", UrlFunc(starttime, endtime), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Auth-Email", "xxxxxxxxxx")
	req.Header.Set("X-Auth-Key", "xxxxxxxxxx")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
	marker = endtime
	fmt.Printf("marker 是 %v\n", marker)

}

var end = Uijmio("2020-03-08T22:00:03Z")

func main() {
	marker = end
	for {
		if marker < time.Now().Unix()-180 {
			GetData(marker, marker+10)
		} else {
			time.Sleep(time.Second * 10)
		}
	}

}
