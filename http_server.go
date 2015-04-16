package main

import (
	"api_server/api_handle"
	"api_server/util"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AddressRequestJson struct {
	Address string `json:"address"`
	Count   string `json:"count"`
}

type Status struct {
	Id              string   `json:"idstr"`
	Attitudes_count int      `json:"attitudes_count"`
	Comments_count  int      `json:"comments_count"`
	Created_at      string   `json:"created_at"`
	Text            string   `json:"text"`
	Original_pic    string   `json:"original_pic"`
	Pic_ids         []string `json:"pic_ids"`
}

type UserWeiboAndComment struct {
	Avatar_hd  string `json:"avatar_hd"`
	Created_at string `json:"created_at"`
	Name       string `json:"name"`
	Id         string `json:"idstr"`
	Status     Status `json:"status"`

	Description     string `json:"description"`
	Gender          string `json:"gender"`
	Followers_count int    `json:"followers_count"`
	Lat             string `json:"lat"`
	Lon             string `json:"lon"`
}

type WeiboSearchMid struct {
	Mid  string `json:"mid"`
	Text string `json:"text"`
}

//type WeiboSearch struct {
//	Statuses []WeiboSearchMid `json:"statuses"`
//}

type WeiboUser struct {
	Id        string `json:"idstr"`
	Name      string `json:"name"`
	Avatar_hd string `json:"avatar_hd"`
}

type WeiboComment struct {
	Created_at string    `json:"created_at"`
	Id         string    `json:"mid"`
	Text       string    `json:"text"`
	User       WeiboUser `json:"user"`
}

var access_token string = "2.004t5RdCHB_LqCd7d61482d5iGDbcD"

func getComment(id string) ([]WeiboComment, error) {
	comments_show_url := fmt.Sprintf("https://api.weibo.com/2/comments/show.json")
	search_place_parameters := url.Values{}
	search_place_parameters.Add("access_token", access_token)
	search_place_parameters.Add("id", id)
	json_string, json_string_err := util.HttpGetHandle(comments_show_url, search_place_parameters)
	if json_string_err != nil {
		fmt.Println("读取评论url失败！")
		return nil, json_string_err
	}

	var weibo_comments []WeiboComment
	json_string_to, _ := json.Marshal(json_string["comments"])

	weibo_comments_err := json.Unmarshal(json_string_to, &weibo_comments)
	if weibo_comments_err != nil {
		fmt.Println("评论json失败！")
		return nil, weibo_comments_err
	}
	return weibo_comments, nil
}

func filterJson(raw_json []byte) ([]byte, error) {
	return_jsons := []map[string]interface{}{}

	var user_weibo_and_comments []UserWeiboAndComment
	fmt.Println(string(raw_json))
	user_weibo_and_comment_err := json.Unmarshal(raw_json, &user_weibo_and_comments)
	if user_weibo_and_comment_err != nil {
		fmt.Println("xxxxxx")
		return nil, errors.New("xxx")
	}

	for _, one := range user_weibo_and_comments {
		return_json := make(map[string]interface{})
		return_json["name"] = one.Name
		return_json["avatar_hd"] = one.Avatar_hd
		return_json["created_at"] = one.Created_at
		return_json["id"] = one.Id
		return_json["description"] = one.Description
		return_json["followers_count"] = one.Followers_count
		return_json["gender"] = one.Gender
		return_json["lat"] = one.Lat
		return_json["lon"] = one.Lon
		weibo_comments, weibo_comments_err := getComment(one.Status.Id)
		if weibo_comments_err != nil {
			fmt.Println("评论hhhh失败！")
			return nil, errors.New("评论hhhh失败！")
		}
		return_json["comments"] = weibo_comments
		return_json["status"] = one.Status
		return_jsons = append(return_jsons, return_json)
	}
	ret, _ := json.Marshal(return_jsons)
	return ret, nil
}

func getInfoByAddr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, read_err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if read_err != nil {
		fmt.Fprintf(w, "%+v\n", "读取失败！")
		return
	}
	var address AddressRequestJson
	json_err := json.Unmarshal(result, &address)
	if json_err != nil {
		fmt.Fprintf(w, "%+v\n", "json解析失败！")
		return
	}

	search_place_url := fmt.Sprintf("https://api.weibo.com/2/place/pois/search.json")
	search_place_parameters := url.Values{}
	search_place_parameters.Add("access_token", access_token)
	search_place_parameters.Add("keyword", address.Address)

	return_json_map, http_err := util.HttpGetHandle(search_place_url, search_place_parameters)
	if http_err != nil {
		fmt.Fprintf(w, "%+v\n", "访问place/pois/search的api出错")
		return
	}
	lon := return_json_map["pois"].([]interface{})[0].(map[string]interface{})["lon"].(string)
	lat := return_json_map["pois"].([]interface{})[0].(map[string]interface{})["lat"].(string)

	nearby_users_url := fmt.Sprintf("https://api.weibo.com/2/place/nearby/users.json")
	nearby_users_parameters := url.Values{}
	nearby_users_parameters.Add("access_token", access_token)
	nearby_users_parameters.Add("lat", lat)
	nearby_users_parameters.Add("long", lon)
	nearby_users_parameters.Add("range", "500")
	nearby_users_parameters.Add("count", address.Count)
	nearby_users_json, nearby_users_error := util.HttpGetHandle(nearby_users_url, nearby_users_parameters)

	if nearby_users_error != nil {
		fmt.Fprintf(w, "%+v\n", "访问nearby/users的api出错")
		return
	}

	user_json, _ := json.Marshal(nearby_users_json["users"])
	flter_user_json, _ := filterJson(user_json)
	fmt.Fprintf(w, "%+v", string(flter_user_json))
}

type TextRequestJson struct {
	Text string `json:"text"`
	Page string `json:"page"`
}

func getInfoByText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, read_err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if read_err != nil {
		fmt.Fprintf(w, "%+v\n", "读取失败！")
		return
	}
	var text_requset TextRequestJson
	json_err := json.Unmarshal(result, &text_requset)
	if json_err != nil {
		fmt.Fprintf(w, "%+v\n", "json解析失败！")
		return
	}
	fmt.Println(text_requset.Page + "----------" + text_requset.Text)

	mids := api_handle.HandleInfoByText(text_requset.Text, text_requset.Page)
	if len(mids) == 0 {
		fmt.Fprintf(w, "%+v", `{"error":"text未找到！"}`)
		return
	}
	return_jsons := []map[string]interface{}{}
	return_json_channel := make(chan map[string]interface{}, 1)
	for _, mid := range mids {
		fmt.Println(mid)
		url_string := "http://weibo.cn/comment/" + mid
		r_url_string := "" + mid
		return_json := make(map[string]interface{})
		go func() {
			ret_html := util.HttpGetClientHandle(url_string)
			uid := api_handle.GetUidFromHtml(ret_html)
			real_url := "http://www.weibo.com/" + uid + "/" + r_url_string
			fmt.Println(real_url)
			raw_html := util.HttpGetPanHandle(real_url)
			weibo_mid := api_handle.GetWeiboMidFromHtml(raw_html)
			fmt.Println("..........." + weibo_mid)

			weibo_show_url := "https://api.weibo.com/2/statuses/show.json"
			weibo_show_parameters := url.Values{}
			weibo_show_parameters.Add("access_token", access_token)
			weibo_show_parameters.Add("id", weibo_mid)

			return_json_map, http_err := util.HttpGetHandle(weibo_show_url, weibo_show_parameters)
			if http_err != nil {
				fmt.Fprintf(w, "%+v\n", "访问place/pois/search的api出错")
				return
			}
			fmt.Println(return_json_map)
			weibo_text := return_json_map["text"].(string)
			//fmt.Println(weibo_text)
			return_json["weibo_id"] = weibo_mid
			return_json["url"] = real_url
			return_json["content"] = weibo_text
			//return_jsons = append(return_jsons, return_json)
			return_json_channel <- return_json

		}()
		fmt.Println("!!!!!!!!!!!!")
		fmt.Println(return_jsons)
	}

	for i := 1; i <= 1; i++ {
		return_json := <-return_json_channel
		return_jsons = append(return_jsons, return_json)
	}
	last_json, _ := json.Marshal(return_jsons)
	fmt.Fprintf(w, "%+v", string(last_json))
}

type UrlRequestJson struct {
	Url string `json:"url"`
}

func getInfoByURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, read_err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if read_err != nil {
		fmt.Fprintf(w, "%+v\n", "读取失败！")
		return
	}
	var url_requset UrlRequestJson
	json_err := json.Unmarshal(result, &url_requset)
	if json_err != nil {
		fmt.Fprintf(w, "%+v\n", "json解析失败！")
		return
	}
	fmt.Println(url_requset.Url)

	raw_html := util.HttpGetPanHandle(url_requset.Url)
	//fmt.Println(string(raw_html))
	weibo_mid := api_handle.GetWeiboMidFromHtml(raw_html)
	fmt.Println("..........." + weibo_mid)

	return_json := map[string]interface{}{}

	comments_show_url := fmt.Sprintf("https://api.weibo.com/2/comments/show.json")
	comments_show_parameters := url.Values{}
	comments_show_parameters.Add("access_token", access_token)
	comments_show_parameters.Add("id", weibo_mid)
	comments_show_json, comments_show_error := util.HttpGetHandle(comments_show_url, comments_show_parameters)
	if comments_show_error != nil {
		fmt.Fprintf(w, "%+v\n", "访问nearby/users的api出错")
		return
	}

	statuses_show_url := fmt.Sprintf("https://api.weibo.com/2/statuses/show.json")
	statuses_show_parameters := url.Values{}
	statuses_show_parameters.Add("access_token", access_token)
	statuses_show_parameters.Add("id", weibo_mid)
	statuses_show_json, statuses_show_error := util.HttpGetHandle(statuses_show_url, statuses_show_parameters)
	if statuses_show_error != nil {
		fmt.Fprintf(w, "%+v\n", "访问nearby/users的api出错")
		return
	}

	return_json["text"] = statuses_show_json["text"]
	return_json["user"] = statuses_show_json["user"]
	return_json["pic_urls"] = statuses_show_json["pic_urls"]
	return_json["thumbnail_pic"] = statuses_show_json["thumbnail_pic"]
	return_json["comments"] = comments_show_json["comments"]
	return_json["source"] = "weibo"

	last_json, _ := json.Marshal(return_json)
	fmt.Fprintf(w, "%+v", string(last_json))
}

func getTextByApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result, read_err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if read_err != nil {
		fmt.Fprintf(w, "%+v\n", "读取失败！")
		return
	}
	var text_requset TextRequestJson
	json_err := json.Unmarshal(result, &text_requset)
	if json_err != nil {
		fmt.Fprintf(w, "%+v\n", "json解析失败！")
		return
	}
	fmt.Println(text_requset.Page + "----------" + text_requset.Text)

	return_jsons := []map[string]interface{}{}
	mids := HandleInfoByApi(text_requset.Text, text_requset.Page)
	fmt.Println(mids)
	for _, mid := range mids {
		return_json := make(map[string]interface{})
		fmt.Println(mid)
		return_json["weibo_id"] = mid.Mid
		return_json["url"] = ""
		return_json["content"] = mid.Text
		return_jsons = append(return_jsons, return_json)
		fmt.Println("!!!!!!!!!!!!")
		fmt.Println(return_jsons)
	}
	last_json, _ := json.Marshal(return_jsons)
	fmt.Fprintf(w, "%+v", string(last_json))
}

func HandleInfoByApi(text string, page string) []WeiboSearchMid {
	search_topics_url := fmt.Sprintf("https://api.weibo.com/2/search/topics.json")
	search_topics_parameters := url.Values{}
	search_topics_parameters.Add("access_token", access_token)
	search_topics_parameters.Add("q", text)
	search_topics_json, search_topics_err := util.HttpGetHandle(search_topics_url, search_topics_parameters)

	if search_topics_err != nil {
		fmt.Println(search_topics_err.Error())
		return nil
	}

	fmt.Println(search_topics_json["statuses"])
	//	return search_topics_json["statuses"].([]string)
	var weibo_search_mids []WeiboSearchMid

	json_string_to, _ := json.Marshal(search_topics_json["statuses"])

	weibo_search_err := json.Unmarshal(json_string_to, &weibo_search_mids)

	//	weibo_search_err := json.Unmarshal(search_topics_json["statuses"], &weibo_search_mids)
	if weibo_search_err != nil {
		fmt.Println(weibo_search_err.Error())
		return nil
	}
	fmt.Println(weibo_search_mids)
	return weibo_search_mids
}

func main() {
	http.HandleFunc("/getInfoByAddr", getInfoByAddr)
	http.HandleFunc("/getInfoByText", getInfoByText)

	http.HandleFunc("/getTextByApi", getTextByApi)

	http.HandleFunc("/getInfoByURL", getInfoByURL)
	http.ListenAndServe(":9091", nil)
}
