package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

var cookies = "SINAGLOBAL=1789977985899.8955.1405930261358; __gads=ID=1e3353c01a843dcd:T=1416881562:S=ALNI_MZxzk8FjZ5cwlzpFFCsTKIAo3roNg; lzstat_uv=29031772931977386616|2893156; YF-V5-G0=55fccf7be1706b6814a78384fa94e30c; _s_tentry=picture.youth.cn; YF-Ugrow-G0=98bd50ad2e44607f8f0afd849e67c645; Apache=7016120641492.308.1424148055909; ULV=1424148055925:31:3:1:7016120641492.308.1424148055909:1423412248574; YF-Page-G0=f70469e0b5607cacf38b47457e34254f; lzstat_ss=1973461882_0_1425423385_2893156; TC-V5-G0=10672b10b3abf31f7349754fca5d2248; TC-Page-G0=b05711a62e11e2c666cc954f2ef362fb; TC-Ugrow-G0=5e22903358df63c5e3fd2c757419b456; appkey=; login_sid_t=d086d3e838f8ca1536c36b506b8db6a9; WBStore=bf5c852d78c6d23c|undefined; ULOGIN_IMG=14278546669585; myuid=2412690033; UOR=,,login.sina.com.cn; SUS=SID-5211538803-1427855453-GZ-5wcvi-0fbbe1459001b626855941e99c59fe48; SUE=es%3Dd976be40e2caefc77e08f1c258bcf3b3%26ev%3Dv1%26es2%3Da65019166ea6f2288fd0436689d44add%26rs0%3DOVJsC65xAVkvpuyl2PYGv%252F7CnyI2Xa7wJRkUQzAYTa6om8uCn3HE1oQew%252BbaNs0UZrDVZhFkJLi2Af5nFFJWDFFAsLvRVSUyMJEX4G0juntYlavZJKLsVuVqd17HlZFVITiTTd6tqlFIKQRa5fqektGnNgHH%252BMexnBO%252Bwx%252FznHQ%253D%26rv%3D0; SUP=cv%3D1%26bt%3D1427855453%26et%3D1427941853%26d%3Dc909%26i%3Dfe48%26us%3D1%26vf%3D0%26vt%3D0%26ac%3D2%26st%3D0%26uid%3D5211538803%26name%3D18515062140%26nick%3Dshanyue2014%26fmp%3D%26lcp%3D; SUB=_2A254HygNDeTxGeNM6lMU8ybEyz-IHXVbbR7FrDV8PUNbuNBeLXHzkW8Rx6WJGNg5MkzkPlygg96AxIVX9Q..; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9W5ChlIrn1VN8zhMEfzGLsYM5JpX5K2t; SUHB=0N2FyB4WB8ahjl; ALF=1459391453; SSOLoginState=1427855453; un=18515062140; wvr=6"

func HttpGetClientHandle(url string) []byte {
	client := &http.Client{}
	req, req_err := http.NewRequest("GET", url, nil)
	if req_err != nil {
		fmt.Println("生成request时出错！")
		return nil
	}
	req.Header.Set("Content-Type", "text/html; charset=utf-8")
	req.Header.Set("Cookie", cookies)
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
	return html
}

func HttpGetPanHandle(url string) []byte {
	client := &http.Client{}
	req, req_err := http.NewRequest("GET", url, nil)
	if req_err != nil {
		fmt.Println("生成request时出错！")
		return nil
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")
	req.Header.Set("Cookie", cookies)
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
	return html
}

func HttpGetHandle(url_string string, parameters url.Values) (map[string]interface{}, error) {
	Url, url_error := url.Parse(url_string)
	if url_error != nil {
		fmt.Println("读取真实url失败！")
		return nil, errors.New("读取真实url失败！")
	}
	Url.RawQuery = parameters.Encode()
	real_api_url := Url.String()
	response, err_response := http.Get(real_api_url)
	if err_response != nil {
		fmt.Println("访问url失败！")
		return nil, err_response
	}
	response_read, response_read_err := ioutil.ReadAll(response.Body)
	if response_read_err != nil {
		fmt.Println("读取url失败！")
		return nil, response_read_err
	}
	var json_data map[string]interface{}
	if json_data_err := json.Unmarshal(response_read, &json_data); json_data_err != nil {
		fmt.Println("解析json失败！")
		return nil, json_data_err
	}
	return json_data, nil
}
