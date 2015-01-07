package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AddressRequestJson struct {
	Address string `json:"address"`
}

type Status struct {
	Attitudes_count int    `json:"attitudes_count"`
	Comments_count  int    `json:"comments_count"`
	Created_at      string `json:"created_at"`
	Text            string `json:"text"`
	Original_pic    string `json:"original_pic"`
}

type UserWeiboAndComment struct {
	Avatar_hd  string `json:"avatar_hd"`
	Created_at string `json:"created_at"`
	Name       string `json:"name"`
	Id         string `json:"id"`
	Status     Status `json:"status"`
}

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

func httpGetHandle(url_string string, parameters url.Values) (map[string]interface{}, error) {
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

//func filterJson(raw_json []map[string]interface{}) (map[string]interface{}, error) {
//fmt.Println(raw_json)
//return nil, nil
//}
func getComment(id string) ([]WeiboComment, error) {
	comments_show_url := fmt.Sprintf("https://api.weibo.com/2/comments/show.json")
	search_place_parameters := url.Values{}
	search_place_parameters.Add("access_token", access_token)
	search_place_parameters.Add("id", id)
	json_string, json_string_err := httpGetHandle(comments_show_url, search_place_parameters)
	if json_string_err != nil {
		fmt.Println("读取评论url失败！")
		return nil, json_string_err
	}
	//fmt.Println(json_string["comments"])
	var weibo_comments []WeiboComment
	json_string_to, _ := json.Marshal(json_string["comments"])
	//fmt.Printf("%+v", string(json_string_to))
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
	user_weibo_and_comment_err := json.Unmarshal(raw_json, &user_weibo_and_comments)
	if user_weibo_and_comment_err != nil {
		fmt.Println("xxxxxx")
		return nil, errors.New("xxx")
	}
	//fmt.Printf("%+v", user_weibo_and_comments)
	for _, one := range user_weibo_and_comments {
		return_json := make(map[string]interface{})
		return_json["Name"] = one.Name
		return_json["Avatar_hd"] = one.Avatar_hd
		return_json["Created_at"] = one.Created_at
		return_json["Id"] = one.Id
		weibo_comments, weibo_comments_err := getComment(one.Id)
		if weibo_comments_err != nil {
			fmt.Println("评论hhhh失败！")
			return nil, errors.New("评论hhhh失败！")
		}
		return_json["Comments"] = weibo_comments
		return_json["Status"] = one.Status
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

	return_json_map, http_err := httpGetHandle(search_place_url, search_place_parameters)
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
	nearby_users_json, nearby_users_error := httpGetHandle(nearby_users_url, nearby_users_parameters)

	if nearby_users_error != nil {
		fmt.Fprintf(w, "%+v\n", "访问nearby/users的api出错")
		return
	}

	//fmt.Println(nearby_users_json["users"])
	//ss := nearby_users_json["users"].([]map[string]interface{})
	//filterJson(ss)
	//for i := range nearby_users_json["user"].([]map[string]interface{}) {
	//fmt.Println(i)
	//fmt.Println("___________________________")
	//}

	user_json, _ := json.Marshal(nearby_users_json["users"])
	//re_json, _ := json.Marshal(nearby_users_json)
	/*
		var user_weibo_and_comments []UserWeiboAndComment
		user_weibo_and_comment_err := json.Unmarshal(re_json, &user_weibo_and_comments)
		if user_weibo_and_comment_err != nil {
			fmt.Println("xxxxxx")
			return
		}
		fmt.Printf("%+v", user_weibo_and_comments)
	*/
	flter_user_json, _ := filterJson(user_json)
	//flter_user_json, _ = json.Marshal(flter_user_json)

	fmt.Fprintf(w, "%+v", string(flter_user_json))
}

func main() {
	http.HandleFunc("/getInfoByAddr", getInfoByAddr)
	http.ListenAndServe(":9091", nil)
	//getComment("3734712990139762")
}
