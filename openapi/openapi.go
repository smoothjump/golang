package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
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

func ShaHmac1(source, secret string) string {
	key := []byte(secret)
	hmac := hmac.New(sha1.New, key)
	hmac.Write([]byte(source))
	signedBytes := hmac.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)
	return signedString
}

func getIso8601Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func formSortedParaStr(para_map *map[string]string) string {
	para_list := make([]string, 0)
	var buffer bytes.Buffer
	for para := range *para_map {
		para_list = append(para_list, para)
	}
	sort.Strings(para_list)
	for _, para := range para_list {
		buffer.WriteString(para)
		buffer.WriteString("=")
		buffer.WriteString((*para_map)[para])
		buffer.WriteString("&")
	}
	buffer.Truncate(buffer.Len() - 1)
	return buffer.String()
}

type AliyunClient struct {
	shared_parameters map[string]string
	endpint           string
}

func (client *AliyunClient) DoGET(customed_parameters map[string]string) {
	client.DoAction(customed_parameters, "GET")
}

func (client *AliyunClient) DoAction(customed_parameters map[string]string, method string) {
	var buffer bytes.Buffer
	buffer.WriteString(method + "&%2F&")
	client.Set("SignatureNonce", getUUID())
	client.Set("Timestamp", getIso8601Now())
	secret := client.Get("access_key") + "&"
	paras := client.shared_parameters
	delete(paras, "AccessSecret")
	for k, v := range customed_parameters {
		paras[k] = v
	}
	fmt.Println(paras)
	buffer.WriteString(percentEncode(formSortedParaStr(&paras)))
	fmt.Println(paras)
	canonicalizedQueryString := buffer.String()
	// fmt.Println(buffer.String())
	signature := ShaHmac1(canonicalizedQueryString, secret)
	// fmt.Println(canonicalizedQueryString)
	// fmt.Println(signature)
	buffer.Truncate(0)
	buffer.WriteString("http://")
	buffer.WriteString(client.endpint)
	buffer.WriteString("/?")
	buffer.WriteString(formSortedParaStr(&paras))
	buffer.WriteString("&Signature=" + signature)
	url := buffer.String()
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		fmt.Println(resp.Body)
	}
	defer resp.Body.Close()
}

func (client *AliyunClient) Get(key string) string {
	if client.shared_parameters == nil {
		client.shared_parameters = make(map[string]string)
	}
	if _, ok := client.shared_parameters[key]; ok {
		return client.shared_parameters[key]
	} else {
		return ""
	}
}

func (client *AliyunClient) Set(key string, value string) {
	if client.shared_parameters == nil {
		client.shared_parameters = make(map[string]string)
	}
	client.shared_parameters[key] = percentEncode(value)
}

func (client *AliyunClient) LoadConfig(product string) {
	config := viper.New()
	config.AddConfigPath("./openapi")
	config.SetConfigName("config")
	config.SetConfigType("yaml")
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	product_conf := config.GetStringMapString("openapi." + product)
	client.Set("AccessKeyId", config.GetString("openapi.AccessKeyId"))
	client.Set("AccessSecret", config.GetString("openapi.AccessSecret"))
	client.Set("RegionId", product_conf["regionid"])
	client.Set("Version", product_conf["version"])
	client.endpint = product_conf["endpoint"]
	client.Set("SignatureMethod", "HMAC-SHA1")
	client.Set("Format", "JSON")
	client.Set("SignatureType", "")
	client.Set("SignatureVersion", "1.0")
}

func main() {
	pp := make(map[string]string)
	pp["Action"] = "DescribeDrdsInstances"
	fmt.Println(getIso8601Now())
	client := AliyunClient{}
	client.LoadConfig("drds")
	client.DoGET(pp)
}
