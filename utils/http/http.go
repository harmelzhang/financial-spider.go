package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// 封装HTTP客户端
type HttpClient struct {
	// HTTP客户端
	client *http.Client
	// 请求地址
	url string
	// HTTP请求头
	headers map[string]string
}

// 创建一个新的HTTP客户端实例
func New(url string, timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
		url:     url,
		headers: make(map[string]string),
	}
}

// 设置默认请求头
func (c *HttpClient) SetHeader(key, value string) {
	c.headers[key] = value
}

// 批量设置请求头
func (c *HttpClient) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		c.headers[k] = v
	}
}

// 发送GET请求
func (c *HttpClient) Get(params map[string]string) ([]byte, *http.Response, error) {
	return c.Request(http.MethodGet, params, nil)
}

// 发送POST请求
func (c *HttpClient) Post(params map[string]string, body any) ([]byte, *http.Response, error) {
	return c.Request(http.MethodPost, params, body)
}

// 通用请求方法
func (c *HttpClient) Request(method string, params map[string]string, body any) ([]byte, *http.Response, error) {
	// 构建完整URL
	fullURL, err := c.buildURL(params)
	if err != nil {
		return nil, nil, fmt.Errorf("build url error, %w", err)
	}

	// 处理请求体
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("json marshal body error, %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	// 创建请求
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("new http request error, %w", err)
	}

	// 设置请求头
	c.setRequestHeaders(req)

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, resp, fmt.Errorf("request url error, %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, resp, fmt.Errorf("request url failed status code is %d", resp.StatusCode)
	}

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, fmt.Errorf("read response body error, %w", err)
	}

	return respBody, resp, nil
}

// 构建完整URL，包含查询参数
func (c *HttpClient) buildURL(params map[string]string) (string, error) {
	// 处理基础URL
	fullPath, err := url.Parse(c.url)
	if err != nil {
		return "", err
	}

	// 添加查询参数
	query := fullPath.Query()
	for k, v := range params {
		query.Add(k, v)
	}
	fullPath.RawQuery = query.Encode()

	return fullPath.String(), nil
}

// 设置请求头
func (c *HttpClient) setRequestHeaders(req *http.Request) {
	// 设置其他请求头
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
}

// ParseResponse 将响应体解析为指定的结构体
func ParseResponse[T any](data []byte) (*T, error) {
	if len(data) == 0 {
		return nil, errors.New("响应体为空")
	}

	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("解析响应体失败: %w", err)
	}

	return &result, nil
}
