/*
* 为了练习用的自带的encoding/json包
*/
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    //"reflect"
    "os/exec"
    "strconv"
)

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

//取推荐列表的参数
type commend struct {
    pagesize int
    page    int
    version  string
}

var commendApi string = "http://mobileapi.5sing.kugou.com/song/getRecommendDailyList"

var getMp3Api string = "http://mobileapi.5sing.kugou.com/song/getSongUrl"

func main() {
    com := commend{10, 1, ""}
    sendUrl := commendApi + "?pagesize=" + strconv.Itoa(com.pagesize) + "&page=" + strconv.Itoa(com.page)
    response, _ := http.Get(sendUrl)
    defer response.Body.Close()
    data, _ := ioutil.ReadAll(response.Body)
    //fmt.Println(string(data))

    if response.StatusCode != 200 {
        fmt.Println("failed")
        return
    }
    list, err := decodeJsonDatas(data)
    if err != nil {
        fmt.Println(err)
        return
    }

    urls, err := getAllMp3Urls(list)
    if err != nil {
        fmt.Println(err)
        return
    }

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
    f,err:=os.Create("mp3.txt")
    if err != nil {
        panic(err)
    }
    for _, v := range urls {
          f.WriteString(v + "\n")
    }
    defer f.Close()
    return nil
}

//解析json数据
func decodeJsonDatas(data []uint8) (myCondItem, error) {

    items := []condToMp3{}
    list := myCondItem{items}

    //解析数据
    var f interface{}
    err := json.Unmarshal(data, &f)
    if err != nil {
        return list, errors.New("json error")
    }
    m := f.(map[string]interface{})

    for _, v := range m {
        switch vv := v.(type) {
        case map[string]interface{}:
            for _, y := range vv {
                switch yy := y.(type) {
                case []interface{}:
                    for _, p := range yy {
                        switch pp := p.(type) {
                        case map[string]interface{}:
                            var st, si string
                            if val, ok := pp["SongType"].(string); ok {
                                st = val

                            }
                            if val, ok := pp["SongId"].(string); ok {
                                si = val

                            }
                            item1 := condToMp3{SongType: st, SongId: si}
                            list.AddItem(item1)
                        }
                    }
                }
            }
        default:
            //fmt.Println("type:", reflect.TypeOf(vv))
        }
    }
    return list, nil
}

//取mp3地址
func getAllMp3Urls(list myCondItem) ([]string, error) {

    var url []string
    //url := make([]string, 0)
    for _, i := range list.Items {
        Mp3Url := getMp3Api + "?songtype=" + i.SongType + "&songid=" + i.SongId

        res, _ := http.Get(Mp3Url)
        defer res.Body.Close()
        song, _ := ioutil.ReadAll(res.Body)
        //解json
        var j interface{}
        err := json.Unmarshal(song, &j)
        if err != nil {
            return url, errors.New("json error")
        }
        sl := j.(map[string]interface{})

        for _, ee := range sl {
            switch ff := ee.(type) {
            case map[string]interface{}:
                var sq string
                if v, ok := ff["squrl"].(string); ok {
                    sq = v

                    url = append(url, sq)
                }
            }
        }
    }
    return url, nil
}
