//package api_handle
package main

import (
	//"bytes/Reader"
	"fmt"
	"net/http"
	//"strings/Reader"
	"net/url"
	//"io"
	//"io/Reader"
	"io/ioutil"
	"strings"
)

func init() {
	fmt.Println("api_handle包完成初始化")
}

func Test(text string) {
	fmt.Println(text)
}

//func httpDo(url string, param Reader) {
func httpDo(url string, param strings.Reader) {

	client := &http.Client{}
	req, req_err := http.NewRequest("POST", url, param)
	if req_err != nil {
		fmt.Println("生成request时出错！")
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//req.Header.Set("Cookie", "name=anny")

	resp, err := client.Do(req)

	defer resp.Body.Close()

	resp, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		fmt.Println("生成request时出错！")
		return
	}

	//fmt.Println(string(resp))
}

func HandleInfoByText(text string) {
	//fmt.Println("for test!!!")
	//url := "http://weibo.cn/search/"
	cookie := "_T_WM=0806963628dfac339b7d53cfa588bddc; gsid_CTandWM=4uro7c471TkstF49jIlZplRL4cX; SUB=_2A255qlyNDeTxGeNM6lMU8ybEyz-IHXVbVWTFrDV6PUJbrdAKLRHEkW0KLIsGFiEEvCDPcMSmoGtuE6Jrmg..; M_WEIBOCN_PARAMS=rl%3D1"
	url_string := "http://14.17.120.252:9091/getInfoByAddr"
	data := url.Values{}
	//data.Set("keyword", text)
	//data.Set("smblog", "搜微博")
	data.Set("address", "东方银座")
	param := strings.NewBufferString(data.Encode())
	httpDo(url, param)
	//fmt.Println(cookie + text + url)
}

func main() {
	read := strings.NewReader("name=邓景")
	fmt.Println(*read)
	HandleInfoByText("过年")
}

/*
func httpDo(url string, param strings.Reader) {
	client := &http.Client{}
	req, req_err := http.NewRequest("POST", url, param)
	if req_err != nil {
		fmt.Println("生成request时出错！")
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	resp, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		fmt.Println("生成request时出错！")
		return
	}
}

func HandleInfoByText(text string) {
	cookie := "_T_WM=0806963628dfac339b7d53cfa588bddc; gsid_CTandWM=4uro7c471TkstF49jIlZplRL4cX; SUB=_2A255qlyNDeTxGeNM6lMU8ybEyz-IHXVbVWTFrDV6PUJbrdAKLRHEkW0KLIsGFiEEvCDPcMSmoGtuE6Jrmg..; M_WEIBOCN_PARAMS=rl%3D1"
	url_string := "http://14.17.120.252:9091/getInfoByAddr"
	data := url.Values{}
	//data.Set("keyword", text)
	//data.Set("smblog", "搜微博")
	data.Set("address", "东方银座")
	param := strings.NewReader(data.Encode())
	httpDo(url, param)
}*/

/*
func HandleInfoByText(text string) {
	client := &http.Client{}
	//url_string := "http://14.17.120.252:9091/getInfoByAddr"
	//cookie := "_T_WM=0806963628dfac339b7d53cfa588bddc; gsid_CTandWM=4uro7c471TkstF49jIlZplRL4cX; SUB=_2A255qlyNDeTxGeNM6lMU8ybEyz-IHXVbVWTFrDV6PUJbrdAKLRHEkW0KLIsGFiEEvCDPcMSmoGtuE6Jrmg..; M_WEIBOCN_PARAMS=rl%3D1"
	cookie := "_T_WM=0806963628dfac339b7d53cfa588bddc; gsid_CTandWM=4uro7c471TkstF49jIlZplRL4cX; SUB=_2A255q_1YDeTxGeNM6lMU8ybEyz-IHXVbV4MQrDV6PUJbrdANLRfFkW2OeRJlMEm5d6-fTL89Bw7JoCgcJg.."
	url_string := "http://weibo.cn/search/"
	//data := url.Values{}
	//data.Set("address", "cctv")
	//raw_json := fmt.Sprintf(`{"keyword": "%s", "smblog": "%s"}`, text, "搜微博")D
	//raw_json := fmt.Sprintf(`{"keyword": "cctv", "smblog": "搜微博"}`)
	raw_json := fmt.Sprintf(`{"keyword": "cctv", "smblog": "%E6%90%9C%E5%BE%AE%E5%8D%9A"}`)
	//raw_json := fmt.Sprintf(`{"address": "%s"}`, text)
	fmt.Println(raw_json)
	jsonStr := []byte(raw_json)

	data := bytes.NewBuffer(jsonStr)

	//param := strings.NewReader(data.Encode())
	fmt.Println(data)

	req, req_err := http.NewRequest("POST", url_string, data)

	//req.AddCookie(req.Cookie(cookie))
	req.Header.Set("Cookie", cookie)

	if req_err != nil {
		fmt.Println("生成request时出错！")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求访问时出错！")
		return
	}

	defer resp.Body.Close()
	json, json_err := ioutil.ReadAll(resp.Body)
	if json_err != nil {
		fmt.Println("读取网页为json时失败！")
		return
	}
	fmt.Println("---------------!!!----------------")
	fmt.Println(string(json))
}
*/
