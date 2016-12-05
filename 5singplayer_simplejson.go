/*
* 用simplejson解析json数据，其它同main.go一样
*/
package main

import (
    "bufio"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    //"reflect"
    "strconv"

    "github.com/bitly/go-simplejson"
)

//取推荐列表的参数
type commend struct {
    pagesize int
    page    int
    version  string
}

//组合成获取mp3条件的struct
type condToMp3 struct {
    SongType, SongId string
}

type myCondItem struct {
    Items []condToMp3

}

func (mycond *myCondItem) AddItem(item condToMp3) []condToMp3 {
    mycond.Items = append(mycond.Items, item)
    return mycond.Items

}

var commendUrl string = "http://mobileapi.5sing.kugou.com/song/getRecommendDailyList"

var getMp3Api string = "http://mobileapi.5sing.kugou.com/song/getSongUrl"

func main() {
    com := commend{10, 1, ""}
    sendUrl := commendUrl + "?pagesize=" + strconv.Itoa(com.pagesize) + "&page=" + strconv.Itoa(com.page)
    response, _ := http.Get(sendUrl)
    defer response.Body.Close()
    data, _ := ioutil.ReadAll(response.Body)

    datastr := string(data)

    jsonData, err := simpleJsonDecode(datastr)
    if err != nil {
        fmt.Println(err)
        return
    }

    list, err := getListData(jsonData)
    if err != nil {
        fmt.Println(err)
        return
    }

    urls, err := getAllMp3Urls(list)
    if err != nil {
        fmt.Println(err)
        return
    }
    //fmt.Println(urls)

    //写入文件
    e := writeToFile(urls)
    if e != nil {
        fmt.Println(err)
    }

    //执行播放
    cmd := exec.Command("mpg123","-@mp3.txt")
    err = cmd.Run()
    if err != nil {
        fmt.Println("cmd.Output: ", err)
        return
    }
	
}

//写入文件
func writeToFile(urls []string) error {

    f, err := os.Create("mp3.txt")
    defer f.Close()

    outWriter := bufio.NewWriter(f)
    for _, v := range urls {
        outWriter.WriteString(v + "\n")
    }
    outWriter.Flush()
    return nil
}

func simpleJsonDecode(datastr string) ([]interface{}, error) {

    var list []interface{}
    var err error
    result, err := simplejson.NewJson([]byte(datastr))

    if err != nil {
        return list, errors.New("json format error")
    }
    if _, ok := result.CheckGet("success"); !ok {
        return list, errors.New("get data failed!")
    }
    list, err = result.Get("data").Get("list").Array()

    if err != nil {
        return list, errors.New("get failed!")
    }
    return list, nil
}

//解析json数据
func getListData(lists []interface{}) (myCondItem, error) {

    items := []condToMp3{}
    list := myCondItem{items}
    for _, v := range lists {
        switch vv := v.(type) {
        case map[string]interface{}:
            var st, si string
            if val, ok := vv["SongType"].(string); ok {
                st = val

            }
            if val, ok := vv["SongId"].(string); ok {
                si = val

            }

            item1 := condToMp3{SongType: st, SongId: si}
            list.AddItem(item1)
        }
    }
    return list, nil
}

//取mp3地址
func getAllMp3Urls(list myCondItem) ([]string, error) {

    var url []string
    for _, i := range list.Items {
        Mp3Url := getMp3Api + "?songtype=" + i.SongType + "&songid=" + i.SongId

        res, _ := http.Get(Mp3Url)
        defer res.Body.Close()
        song, _ := ioutil.ReadAll(res.Body)
        //解json
        result, err := simplejson.NewJson([]byte(song))
        if err != nil {
            return url, errors.New("json format error")
        }
        if _, ok := result.CheckGet("success"); !ok {
            return url, errors.New("get data failed!")
        }
        rec, _ := result.Get("data").Map()
        if v, ok := rec["squrl"].(string); ok {
            url = append(url, v)
        }
    }
    return url, nil
}
