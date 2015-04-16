package api_handle

//package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	//"reflect"
	"regexp"
)

func init() {
	fmt.Println("handle_info_text初始化")
}

func Test(text string) {
	fmt.Println(text)
}

func filterMidString(src string) string {
	re, _ := regexp.Compile("M_(.*)")
	last_src := re.FindStringSubmatch(src)
	return last_src[1]
}

func HandleInfoByText(text string, page string) []string {
	client := &http.Client{}
	cookie := "_T_WM=e79bbdc428bd56d824d340826338c489; WEIBOCN_WM=5091_0010; SUB=_2A254Hyj-DeTxGeNM6lMU8ybEyz-IHXVb4Ei2rDV6PUJbvNBeLUvGkW15C18vqYeSvZcvsQjyxKHiRBtMHw..; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9W5ChlIrn1VN8zhMEfzGLsYM5JpX5K-t; SUHB=0N2FyB4WB8ahjl; gsid_CTandWM=4uPhe7771jwoPMy0m6dDclRL4cX; SUE=es%3D33d4ef5a34cecd665f81f3b3d5e33ffb%26ev%3Dv1%26es2%3D04bc4aa84a2d1961ed1786fd3836c86f%26rs0%3DZIBo%252FrtP6JFYPKqsSVHVLFl6a8NcU1uVf3v9erbUhVjBneeo%252BF7lVeqN5HAEJF5k8lCpDHuY%252Fst3%252BEkmQWtUvlz%252B5s%252BPmLIFEQaEHiwHD1%252BwPJ42i4phbwk%252BpveFcAyN9JYQzIDKwjbbU3f%252FexsoSuV1pRDI06VgERy5kowJsnw%253D%26rv%3D0; SUP=cv%3D1%26bt%3D1427855535%26et%3D1427941935%26d%3Dc65e%26i%3D6e6f%26us%3D1%26vf%3D%26vt%3D%26ac%3D%26st%3D0%26lt%3D1%26uid%3D5211538803%26user%3D18515062140%26ag%3D1%26name%3D18515062140%26nick%3Dshanyue2014%26sex%3D%26ps%3D0%26email%3D%26dob%3D%26ln%3D%26os%3D%26fmp%3D%26lcp%3D; SUS=SID-5211538803-1427855535-GZ-ds7t1-56dca0ff14a4e55996a0ce36395f6e6f; SSOLoginState=1427855535"

	raw_url_string := fmt.Sprintf("http://weibo.cn/search/mblog?keyword=%s&sort=hot&page=%s", text, page)
	parse_url, _ := url.Parse(raw_url_string)

	param := parse_url.Query().Encode()
	fmt.Println(param)

	real_url_string := "http://weibo.cn/search/mblog?" + param

	fmt.Println(parse_url.Query().Encode())
	fmt.Println(real_url_string)

	req, req_err := http.NewRequest("GET", real_url_string, nil)
	req.Header.Set("Cookie", cookie)

	if req_err != nil {
		fmt.Println("生成request时出错！")
		return nil
	}
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求访问时出错！")
		return nil
	}

	defer resp.Body.Close()

	html, html_err := ioutil.ReadAll(resp.Body)

	if html_err != nil {
		fmt.Println("读取网页为json时失败！")
		return nil
	}
	fmt.Println("---------------!!!----------------")

	doc, doc_err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if doc_err != nil {
		fmt.Print("goquery出错")
		return nil
	}

	ret_strings := []string{}
	doc.Find("[id^='M_']").Each(func(i int, s *goquery.Selection) {
		band, _ := s.Attr("id")
		//fmt.Println(band, reflect.TypeOf(band))
		band = filterMidString(band)
		ret_strings = append(ret_strings, band)
	})
	fmt.Println(ret_strings)
	return ret_strings
}

func filterUidString(raw_uid string) string {
	fmt.Println(raw_uid)
	re, _ := regexp.Compile(".*uid=(.*)&rl.*")
	uid := re.FindStringSubmatch(raw_uid)
	fmt.Println(uid[1])
	return uid[1]
}

func GetUidFromHtml(html []byte) string {
	doc, doc_err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if doc_err != nil {
		fmt.Print("goquery出错")
		return "goquery出错"
	}
	raw_uid, _ := doc.Find("[name='srcuid']").Attr("value")
	fmt.Println(raw_uid)
	return raw_uid
}

func filterWeiboMidString(raw_uid string) string {
	re, _ := regexp.Compile(`mblog&act=([\d]+)`)
	uid := re.FindStringSubmatch(raw_uid)
	return uid[1]
}

func GetWeiboMidFromHtml(html []byte) string {
	fmt.Println("FGHJKKNKNLM..................")
	//	fmt.Println(string(html))

	raw_mid := filterWeiboMidString(string(html))
	fmt.Println("QQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQQ")
	fmt.Println(raw_mid)
	return raw_mid
}

func main() {
	fmt.Println(filterWeiboMidString("weibo.com/a/vpaint/rec?do=mblog&act=3796039976449655 sdsdc"))
}
