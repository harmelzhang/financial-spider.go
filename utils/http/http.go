package http

import (
	"financial-spider.go/config"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Get 请求网络资源，HTTP : GET
func Get(url string) []byte {
	client := &http.Client{Timeout: config.HttpTimeout * time.Second}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Accept", config.HttpAccept)
	req.Header.Add("User-Agent", config.UserAgent[rand.Intn(len(config.UserAgent))])
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("执行网络请求出错 : %s > %s", url, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		log.Fatalf("网络请求出错，Status Code: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("执行网络请求出错 : %s", err)
	}
	return bytes
}
