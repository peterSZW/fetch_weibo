package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"encoding/json"
	"time"
	"math"
)

type FetchError struct {
	ErrorMessage string
}

func (e *FetchError) Error() string {
	return "Fetch data error: " + e.ErrorMessage
}

func get_url(url_path string, file_path string) (map[string]interface{}, bool, error) {
	var (
		client = &http.Client{}
		base_url string = "http://api.weibo.com"
		appkey string = os.Getenv("APPKEY")
		cookies string = os.Getenv("COOKIES")
		overload = false
	)

	url_path = base_url + url_path + "&source=" + appkey;
	req, _ := http.NewRequest("GET", url_path, nil)
	req.Header.Set("Cookie", cookies)

	res, err := client.Do(req)
	if (err != nil) {
		return nil, overload, err
	} else {
		fmt.Printf("%s GET %s\n", res.Status, url_path)
		body, err := ioutil.ReadAll(res.Body)
		if (err != nil) {
			return nil, overload, err
		} else {
			var f interface{}
			json.Unmarshal(body, &f)
			ret := f.(map[string]interface{})
			if (ret["error_code"] != nil && math.Abs(ret["error_code"].(float64) - 10023.0) < 0.1 ) {
				overload = true
			} else if (ret["error_code"] != nil) {
				return nil, overload, &FetchError{ret["error"].(string)}
			}
			ioutil.WriteFile(file_path, body, 0644)
			return ret, overload, nil
		}
	}
	defer res.Body.Close()

	return nil,overload, nil
}

func get_friend(uid string) ([]string, bool, error) {
	item, overload, err := get_url(
		"/2/friendships/friends/bilateral/ids.json?uid=" + uid,
		"friend_go/" + uid,
	)

	if (err != nil || overload) {
		return nil, overload, err
	}
	ids := make([]string, 0, 100)
	array := item["ids"].([]interface{})
	for _, k := range array {
		ids = append(ids, fmt.Sprintf("%.0f", k))
	}
	return ids, overload, nil
}

func get_timeline(uid string) (bool, error) {
	_, overload, err := get_url(
		"/2/statuses/user_timeline.json?uid=" + uid,
		"timeline_go/" + uid,
	)

	return overload, err
}

func user_exist(uid string) bool {
	var path = "friend_go/" + uid
	_, err := os.Stat(path)
	if err == nil { return true}
	if os.IsNotExist(err) {return false}
	return false
}

// get all data from weibo
// errorTimes: after how many errors the function exit?
// errorWait: how long to wait when ocurs error? (in minutes)
// overload: how long to wait when request is out of rate limit? (in minutes)
// total: how many users we will fetch?
func get_all(errorTimes int, errorWait int, overloadWait int, total int) error {
	q := make([]string, 0, 1000)
	errCount := 0
	num := 0

	q = append(q, "1967956383")

	for ;len(q) > 0 && num < total; {
		id := q[0]
		get_timeline(id)
		users, overload, err := get_friend(id)
		if (err != nil) {
			if (errCount > errorTimes) {
				fmt.Printf("error over %d times, exit!\n", errorTimes)
				return err
			} else {
				errCount++
				fmt.Println(err)
				fmt.Printf("error ocurs, retry after %d minute...\n", errorWait)
				time.Sleep(time.Duration(errorWait) * time.Minute)
			}
		} else {
			if (overload) {
				fmt.Printf("request out of rate limit, retry after %d minutes...\n", overloadWait)
				time.Sleep(time.Duration(overloadWait) * time.Minute)
			} else {
				q = q[1:]
				num++
				for j, u := range users {
					if (j > 10) {
						break
					}
					if (!user_exist(u)) {
						q = append(q, u)
					}
				}
			}
		}
	}

	defer fmt.Printf("Get users complete, %d users in total\n", num)

	return nil
}

func main() {
	get_all(50, 1, 10, 10000)
}
