package impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// 发送GET请求
// url:请求地址
// response:请求返回的内容
func Get(url string) (response string) {
	client := http.Client{Timeout: 30 * time.Second}
	resp, error := client.Get(url)
	if error != nil {
		panic(error)
	}
	resp.Header.Add("content-type", "application/json")
	resp.Header.Add("content-type", "charset=UTF-8")
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			// fmt.Println("关闭conn", err)
			break
		} else if err != nil {
			panic(err)
		}
	}

	response = result.String()
	return
}

// 发送POST请求
// url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
// content:请求放回的内容
func Post(url string, data interface{}, contentType string) (content string) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	req.Header.Add("content-type", "charset=UTF-8")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 15 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}

//post
func HttpPost(uri string, data string) ([]byte, error) {
	response, err := http.Post(uri, "application/x-www-form-urlencoded;charset=utf-8", bytes.NewBuffer([]byte(data)))
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

// //HTTPGet get 请求
// func HTTPGet(uri string) ([]byte, error) {
// 	response, err := http.Get(uri)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer response.Body.Close()
// 	if response.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
// 	}
// 	return ioutil.ReadAll(response.Body)
// }

// //post
// func HttpPost(uri string, data string) ([]byte, error) {
// 	response, err := http.Post(uri, "application/x-www-form-urlencoded;charset=utf-8", bytes.NewBuffer([]byte(data)))
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	defer response.Body.Close()

// 	if response.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
// 	}
// 	return ioutil.ReadAll(response.Body)
// }

// //PostJSON post json 数据请求
// func PostJSON(uri string, obj interface{}) ([]byte, error) {
// 	jsonData, err := json.Marshal(obj)
// 	if err != nil {
// 		return nil, err
// 	}
// 	body := bytes.NewBuffer(jsonData)
// 	response, err := http.Post(uri, "application/json;charset=utf-8", body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer response.Body.Close()

// 	if response.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
// 	}
// 	return ioutil.ReadAll(response.Body)
// }

// //PostFile 上传文件
// func PostFile(fieldname, filename, uri string) ([]byte, error) {
// 	fields := []MultipartFormField{
// 		{
// 			IsFile:    true,
// 			Fieldname: fieldname,
// 			Filename:  filename,
// 		},
// 	}
// 	return PostMultipartForm(fields, uri)
// }

// //MultipartFormField 保存文件或其他字段信息
// type MultipartFormField struct {
// 	IsFile    bool
// 	Fieldname string
// 	Value     []byte
// 	Filename  string
// }

// //PostMultipartForm 上传文件或其他多个字段
// func PostMultipartForm(fields []MultipartFormField, uri string) (respBody []byte, err error) {
// 	bodyBuf := &bytes.Buffer{}
// 	bodyWriter := multipart.NewWriter(bodyBuf)

// 	for _, field := range fields {
// 		if field.IsFile {
// 			fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
// 			if e != nil {
// 				err = fmt.Errorf("error writing to buffer , err=%v", e)
// 				return
// 			}

// 			fh, e := os.Open(field.Filename)
// 			if e != nil {
// 				err = fmt.Errorf("error opening file , err=%v", e)
// 				return
// 			}
// 			defer fh.Close()

// 			if _, err = io.Copy(fileWriter, fh); err != nil {
// 				return
// 			}
// 		} else {
// 			partWriter, e := bodyWriter.CreateFormField(field.Fieldname)
// 			if e != nil {
// 				err = e
// 				return
// 			}
// 			valueReader := bytes.NewReader(field.Value)
// 			if _, err = io.Copy(partWriter, valueReader); err != nil {
// 				return
// 			}
// 		}
// 	}

// 	contentType := bodyWriter.FormDataContentType()
// 	bodyWriter.Close()

// 	resp, e := http.Post(uri, contentType, bodyBuf)
// 	if e != nil {
// 		err = e
// 		return
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, err
// 	}
// 	respBody, err = ioutil.ReadAll(resp.Body)
// 	return
// }
