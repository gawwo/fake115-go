package utils

import (
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"github.com/gawwo/fake115-go/config"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 不要反复创建Client，它应该全局唯一
func NewClient() *http.Client {
	dialer := &net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 15 * time.Second,
	}
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 2 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		DialContext:         dialer.DialContext,
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 50,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Minute,
	}
	return client
}

var client = NewClient()

func Request(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	// 设置可能缺少的默认参数
	if _, ok := headers["Connection"]; !ok {
		headers["Connection"] = "keep-alive"
	}
	if _, ok := headers["Accept"]; !ok {
		headers["Accept"] = "*/*"
	}
	if _, ok := headers["Accept-Encoding"]; !ok {
		headers["Accept-Encoding"] = "gzip, deflate"
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	var (
		res          *http.Response
		requestError error
	)

	for i := 1; ; i++ {
		res, requestError = client.Do(req)
		if requestError == nil && res.StatusCode < 500 {
			break
		}

		var errMsg string

		// 网络可能有问题
		if requestError != nil {
			errMsg := fmt.Sprintf("request error: %v", requestError)
			config.Logger.Error(errMsg)
		} else {
			// 对方返回的http status code有问题
			errMsg := fmt.Sprintf("%s request error: HTTP %d", url, res.StatusCode)
			config.Logger.Error(errMsg)
		}

		// 重试到最大次数
		if i >= config.RetryTimes {
			return nil, fmt.Errorf(errMsg)
		}

		time.Sleep(time.Second)
	}
	return res, nil
}

func get(urlGet string, headers map[string]string, data map[string]string, withResponse bool) (
	[]byte, *http.Response, error) {
	// 尝试拼接get的参数
	if data != nil {
		getData := url.Values{}
		for k, v := range data {
			getData.Set(k, v)
		}
		urlGet = urlGet + "?" + getData.Encode()
	}

	if headers == nil {
		headers = map[string]string{}
	}

	res, err := Request(http.MethodGet, urlGet, nil, headers)
	if err != nil {
		return nil, nil, err
	}
	return readBody(res, withResponse)
}

func Get(urlGet string, headers map[string]string, data map[string]string) ([]byte, error) {
	body, _, err := get(urlGet, headers, data, false)
	return body, err
}

func GetResponse(urlGet string, headers map[string]string, data map[string]string) ([]byte, *http.Response, error) {
	return get(urlGet, headers, data, true)
}

func PostForm(urlPost string, headers map[string]string, data map[string]string) ([]byte, error) {
	postData := url.Values{}
	for k, v := range data {
		postData.Set(k, v)
	}

	// 不设置content-type，对方就可能认为没发送form body
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	}

	body, _, err := postByte(urlPost, headers, postData, false)
	return body, err
}

func readBody(res *http.Response, withResponse bool) ([]byte, *http.Response, error) {
	var reader io.ReadCloser
	defer res.Body.Close()

	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
	case "deflate":
		reader = flate.NewReader(res.Body)
	default:
		reader = res.Body
	}
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil && err != io.EOF {
		return nil, nil, err
	}

	if withResponse {
		return body, res, nil
	} else {
		return body, nil, nil
	}

}

func postByte(url string, headers map[string]string, data url.Values, withResponse bool) ([]byte, *http.Response, error) {
	if headers == nil {
		headers = map[string]string{}
	}

	dataString := data.Encode()
	if _, ok := headers["Content-Length"]; !ok {
		headers["Content-Length"] = strconv.Itoa(len(dataString))
	}

	res, err := Request(http.MethodPost, url, strings.NewReader(dataString), headers)
	if err != nil {
		return nil, nil, err
	}

	return readBody(res, withResponse)
}
