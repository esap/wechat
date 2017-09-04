package util

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

// GetJson 发送GET请求解析json
func GetJson(uri string, v interface{}) error {
	r, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// GetXml 发送GET请求并解析xml
func GetXml(uri string, v interface{}) error {
	r, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return xml.NewDecoder(r.Body).Decode(v)
}

// GetBody 发送GET请求，返回body字节
func GetBody(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get err: uri=%v , statusCode=%v", uri, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

//PostJson 发送Json格式的POST请求
func PostJson(uri string, obj interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(uri, "application/json;charset=utf-8", buf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

//PostJsonPtr 发送Json格式的POST请求并解析结果到result指针
func PostJsonPtr(uri string, obj interface{}, result interface{}) (err error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(obj)
	if err != nil {
		return
	}

	// debug
	//	println("post body:", buf.String())

	resp, err := http.Post(uri, "application/json;charset=utf-8", buf)
	if err != nil {
		return err
	}

	// debug
	//	bd, _ := ioutil.ReadAll(resp.Body)
	//	println("resp body:", string(bd))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http post error : uri=%v , statusCode=%v", uri, resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(result)
}

// PostFile 上传文件
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

// GetFile 下载文件
func GetFile(filename, uri string) error {
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	return err
}

// MultipartFormField 文件或其他表单数据
type MultipartFormField struct {
	IsFile    bool
	Fieldname string
	Value     []byte
	Filename  string
}

// PostMultipartForm 上传文件或其他表单数据
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
	return ioutil.ReadAll(resp.Body)
}
