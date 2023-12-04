/**
 * @Author: Hardews
 * @Date: 2023/4/5 0:32
 * @Description:
**/

package tool

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func sendPostByJsonForm(url string, text string) ([]byte, error) {
	var reqForm = struct {
		Text string `json:"text"`
	}{text}

	jsonByte, err := json.Marshal(&reqForm)
	if err != nil {
		return nil, err
	}

	var (
		req    *http.Request
		res    *http.Response
		client = &http.Client{}
	)

	req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonByte))
	if err != nil {
		return nil, err
	}

	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	return body, err
}

func sendPostByWWW(Url string, Form map[string]string) ([]byte, error) {
	var (
		err    error
		req    *http.Request
		res    *http.Response
		client = &http.Client{}
	)

	postData := url.Values{}
	for key, val := range Form {
		postData.Add(key, val)
	}

	req, err = http.NewRequest(http.MethodPost, Url, strings.NewReader(postData.Encode()))
	if err != nil {
		return nil, err
	}

	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	return body, err
}
