package main

import (
        "flag"
        "fmt"
        "io"
        "log"
        "net/http"
        "os"
        "strconv"
        "time"
)

//设置游标
var marker int64
var (
        zoneID  string
        email   string
        api_key string
)

func Uijmio(str string) int64 {
        p, _ := time.Parse("2006-01-02T15:04:05Z", str)
        io := p.Unix()
        return io
}

func UrlFunc(starttime, endtime int64) string {
        ApiUrl := "https://api.cloudflare.com/client/v4/zones/" + zoneID + "/logs/received?"
        ApiUrl += "start=" + strconv.FormatInt(starttime, 10) + "&" + "end=" + strconv.FormatInt(endtime, 10)
        return ApiUrl //+ "&count=10"
}

//获取数据
func GetData(starttime, endtime int64) {
        client := &http.Client{}
        req, err := http.NewRequest("GET", UrlFunc(starttime, endtime), nil)
        if err != nil {
                log.Fatal(err)
        }
        //        req.Header.Set("accept-encoding", "gzip")
        req.Header.Set("X-Auth-Email", email)
        req.Header.Set("X-Auth-Key", api_key)
        resp, err := client.Do(req)
        if err != nil {
                log.Fatal(err)
        }
        //bodyText, err := ioutil.ReadAll(resp.Body)
        //if err != nil {
        //      log.Fatal(err)
        //}
        //content := fmt.Sprintf("%s\n", bodyText)

        f, err := os.Create("chimpone_log_" + time.Now().Format("20060102150405") + ".log")
        if err != nil {
                panic(err)
        }
        defer func() { _ = f.Close() }()

        _, err = io.Copy(f, resp.Body)
        if err != nil {
                panic(err)
        }

        //Wrlog(content)
        marker = endtime
        fmt.Printf("marker 是 %v\n", marker)

}

var initTime int64

func Wrlog(content string) {

        f, err := os.OpenFile("a.log", os.O_WRONLY, 0644)
        if err != nil {
                // 打开文件失败处理

        } else {
                // 查找文件末尾的偏移量
                n, _ := f.Seek(0, 2)

                // 从末尾的偏移量开始写入内容
                _, err = f.WriteAt([]byte(content), n)
        }

        defer f.Close()
}

func main() {
        flag.StringVar(&email, "e", "", "email, Default is nil")
        flag.StringVar(&api_key, "a", "", "api_key, Default is nil")
        flag.StringVar(&zoneID, "z", "default value", "zoneID, The zoneID can found in cloudflare overview page, and must be a string type like aabbcc112233!, default it's chimpone.com's zoneID if you not set it ")
        flag.Int64Var(&initTime, "t", time.Now().Unix()-360, "initTime, The start time of you want to search(Timestamp),default it's rightnow if you not set it!")
        flag.Parse()
        marker = initTime
        for {
                if marker < time.Now().Unix()-360 {
                        GetData(marker, marker+300)
                } else {
                        time.Sleep(time.Second * 300)
                        fmt.Println("marker > now - 360 sec\n", marker)
                        fmt.Println(time.Now().Unix() - 360)
                }
        }

}
