package main

import (
    "bytes"
    "fmt"
    "net/url"
    "sort"
    "strings"
    "time"
    "github.com/spf13/viper"
    uuid "github.com/satori/go.uuid"
)

const (
    config_url := ""
)

func getUUID() string {
    id := uuid.NewV4()
    return id.String()
}

func percentEncode(value string) string {
    if value != "" {
        return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(url.QueryEscape(value), "+", "%20"), "*", "%2A"), "%7E", "~")
    } else {
        return ""
    }
}

func getIso8601Now() string {
    return time.Now().UTC().Format(time.RFC3339)
}

func formSortedParaStr(para_map *map[string]string) string {
    para_list := make([]string, 0)
    var buffer bytes.Buffer
    for para, _ := range *para_map {
        para_list = append(para_list, para)
    }
    sort.Strings(para_list)
    for _, para := range para_list {
        buffer.WriteString(para)
        buffer.WriteString("=")
        buffer.WriteString((*para_map)[para])
    }
    return buffer.String()
}

type AliyunClient struct {
    parameters map[string]string
}

func (client *AliyunClient) DoAction() string {

}

func (client *AliyunClient) LoadConfig() string {

}

func main() {
    fmt.Println(percentEncode("+dsdaf+dsdsfa+fafasf+dad+daf"))
    fmt.Println(getIso8601Now())
}

