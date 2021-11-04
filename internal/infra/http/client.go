package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"tiktok-uploader/log"
	"unsafe"
)

type Param map[string]string
type FormatData map[string]interface{}

func getClient() *http.Client {
	//client := &http.Client{}
	return http.DefaultClient
}

func DoGet(uri string, params Param) (*[]byte, error) {
	param := url.Values{}
	Url, err := url.Parse(uri)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for key, value := range params {
		param.Set(key, value)
	}

	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = param.Encode()
	resp, err := getClient().Get(Url.String())
	return execute(resp, err)
}

func DoPost(uri string, params FormatData) (*[]byte, error) {
	bytesData, err := json.Marshal(params)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	req, err := http.NewRequest("POST", uri, reader)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := getClient().Do(req)
	return execute(resp, err)
}


func UploadData(url, name, fileName string, bts *[]byte) (*[]byte, error) {
	reader := bytes.NewReader(*bts)
	return UploadFile(url, name, fileName, reader)
}

func UploadFile(url, name, fileName string, file io.Reader) (*[]byte, error) {
	var bufReader bytes.Buffer

	mpWriter := multipart.NewWriter(&bufReader)
	//字段名必须为media
	writer, err := mpWriter.CreateFormFile(name, fileName)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	io.Copy(writer, file)
	//关闭了才能把内容写入
	mpWriter.Close()

	req, err := http.NewRequest("POST", url, &bufReader)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//从mpWriter中获取content-Type
	req.Header.Set("Content-Type", mpWriter.FormDataContentType())
	resp, err := getClient().Do(req)
	return execute(resp, err)
}

func execute(resp *http.Response, err error) (*[]byte, error) {
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer Close(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		errMsg := (*string) (unsafe.Pointer(&body))
		log.Error(*errMsg)
		return nil, Error(resp.StatusCode, *errMsg)
	}
	return &body, nil
}