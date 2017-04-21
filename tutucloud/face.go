package tutucloud

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// API 接口地址
const API_URL = "https://api.tutucloud.com/v1/face/"

type Keys struct {
	Api_key string
	Api_secret string
}

type FaceApi struct {
	Keys Keys
}

// API 返回 JSON 结构
type Result struct {
	Ret     int                    `json: "ret"`     // 返回码
	Message string                 `json: "message"` // 返回码说明
	Data    map[string]interface{} `json: "data"`    // 数据
	Ttp     int                    `json: "ttp"`     // 服务器时间戳
}

// face API 通用请求方法, 传入接口方法和参数,  返回 JSON map
func (f *FaceApi) Request(method string, image map[string]string, params map[string]string) (*Result, error) {
	url := API_URL + method

	// 通过 FaceApi URL 或 File 获取图片参数
	if image["url"] != "" {
		params["img"] = image["url"]
	} else if image["file"] == "" {
		return nil, errors.New("img is required")
	}

	// 公有 key
	params["api_key"] = f.Keys.Api_key
	// 时间戳
	params["t"] = strconv.Itoa(int(time.Now().Unix()))
	// 参数签名
	params["sign"] = signature(params, f.Keys.Api_secret)

	// 返回 POST 请求实例
	request, err := f.post(url, image, params)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	// if resp.StatusCode != 200 {
	// 	return nil, errors.New("response StatusCode :" + strconv.Itoa(resp.StatusCode))
	// }

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// 解析 JSON
	r := &Result{}
	json.Unmarshal(body, &r)
	return r, nil
}

// HTTP POST
func (f *FaceApi) post(url string, image map[string]string, params map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 图片文件字段
	if image["file"] != "" {
		file, err := os.Open(image["file"])
		defer file.Close()
		if err != nil {
			return nil, err
		}
		fileContents, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormFile("img", filepath.Base(image["file"]))
		if err != nil {
			return nil, err
		}
		part.Write(fileContents)
	}

	// form 参数
	for k, v := range params {
		writer.WriteField(k, v)
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", url, body)	
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, err
}
