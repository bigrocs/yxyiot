package util

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/clbanning/mxj"
	"golang.org/x/crypto/pkcs12"
)

//HTTPGet get 请求
func HTTPGet(uri string) ([]byte, error) {
	response, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//PostForm form  数据请求
func PostForm(url string, obj string) ([]byte, error) {
	reader := strings.NewReader(obj)
	response, err := http.Post(url, "application/x-www-form-urlencoded", reader)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : url=%v , statusCode=%v", url, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//PostJSON post json 数据请求
func PostJSON(uri string, obj interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(uri, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http post error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//PostJSONHeaders post json 数据请求
func PostJSONHeader(uri string, obj interface{}, headers map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer(jsonData)
	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
	response, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

// PostJSONWithRespContentType post json数据请求，且返回数据类型
func PostJSONWithRespContentType(uri string, obj interface{}) ([]byte, string, error) {
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return nil, "", err
	}

	jsonData = bytes.Replace(jsonData, []byte("\\u003c"), []byte("<"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u003e"), []byte(">"), -1)
	jsonData = bytes.Replace(jsonData, []byte("\\u0026"), []byte("&"), -1)

	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(uri, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	contentType := response.Header.Get("Content-Type")
	return responseData, contentType, err
}

//PostFile 上传文件
func PostFile(fieldname, filename, uri string) ([]byte, error) {
	fields := []MultipartFormField{
		{
			IsFile:    true,
			Fieldname: fieldname,
			Filename:  filename,
		},
	}
	return PostMultipartForm(fields, uri)
}

//MultipartFormField 保存文件或其他字段信息
type MultipartFormField struct {
	IsFile    bool
	Fieldname string
	Value     []byte
	Filename  string
}

//PostMultipartForm 上传文件或其他多个字段
func PostMultipartForm(fields []MultipartFormField, uri string) (respBody []byte, err error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for _, field := range fields {
		if field.IsFile {
			fileWriter, e := bodyWriter.CreateFormFile(field.Fieldname, field.Filename)
			if e != nil {
				err = fmt.Errorf("error writing to buffer , err=%v", e)
				return
			}

			fh, e := os.Open(field.Filename)
			if e != nil {
				err = fmt.Errorf("error opening file , err=%v", e)
				return
			}
			defer fh.Close()

			if _, err = io.Copy(fileWriter, fh); err != nil {
				return
			}
		} else {
			partWriter, e := bodyWriter.CreateFormField(field.Fieldname)
			if e != nil {
				err = e
				return
			}
			valueReader := bytes.NewReader(field.Value)
			if _, err = io.Copy(partWriter, valueReader); err != nil {
				return
			}
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, e := http.Post(uri, contentType, bodyBuf)
	if e != nil {
		err = e
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	respBody, err = ioutil.ReadAll(resp.Body)
	return
}

//PostXML perform a HTTP/POST request with XML body
func PostXML(uri string, obj map[string]interface{}) ([]byte, error) {
	mv := mxj.Map(obj)
	xmlData, err := mv.Xml()
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(xmlData)
	response, err := http.Post(uri, "application/xml;charset=utf-8", body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http code error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

//httpWithTLS CA证书
func httpWithTLS(rootCa string, certData string, password string) (client *http.Client, err error) {
	var cert tls.Certificate
	var certD []byte
	if rootCa != "" {
		certD, err = ioutil.ReadFile(rootCa)
		if err != nil {
			return nil, fmt.Errorf("unable to find cert path=%s, error=%v", rootCa, err)
		}
	}
	if certData != "" {
		certD, err = base64.StdEncoding.DecodeString(certData)
		if err != nil {
			return nil, fmt.Errorf("certData 商家秘钥文件转码错误", err)
		}
	}
	cert, err = pkcs12ToPem(certD, password)
	if err != nil {
		return client, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tr := &http.Transport{
		TLSClientConfig:    config,
		DisableCompression: true,
	}

	client = &http.Client{Transport: tr}
	return client, err
}

//pkcs12ToPem 将Pkcs12转成Pem
func pkcs12ToPem(p12 []byte, password string) (cert tls.Certificate, err error) {
	blocks, err := pkcs12.ToPEM(p12, password)
	defer func() {
		if x := recover(); x != nil {
			log.Print(x)
		}
	}()
	if err != nil {
		return cert, err
	}
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	cert, err = tls.X509KeyPair(pemData, pemData)

	if err != nil {
		return cert, err
	}
	return cert, err
}

//PostXMLWithTLS perform a HTTP/POST request with XML body and TLS
// func PostXMLWithTLS(uri string, obj interface{}, ca, key string) ([]byte, error) {
// 	xmlData, err := xml.Marshal(obj)
// 	if err != nil {
// 		return nil, err
// 	}

// 	body := bytes.NewBuffer(xmlData)
// 	client, err := httpWithTLS(ca, key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	response, err := client.Post(uri, "application/xml;charset=utf-8", body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer response.Body.Close()

// 	if response.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("http code error : uri=%v , statusCode=%v", uri, response.StatusCode)
// 	}
// 	return ioutil.ReadAll(response.Body)
// }

//PostJSONHeaders post json 数据请求
func PostJSONWithTLS(uri string, bodyAES string, headers map[string]interface{}, ca string, certData string, password string) ([]byte, map[string][]string, error) {
	body := bytes.NewBuffer([]byte(bodyAES))
	client, err := httpWithTLS(ca, certData, password)
	if err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	for k, v := range headers {
		req.Header.Set(k, v.(string))
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	return res, response.Header, err
}
