package http_util

import (
	"net/http"
	"time"
)

var DefaultHttpClient = &http.Client{
	Timeout: time.Second * 3,
}
/**
带重试的Get请求，最多尝试 1 + retryCnt 次请求，最少尝试1次请求
 */
func GetWithRetry(url string, retryCnt int, header ...http.Header) (*http.Response, error) {

	if retryCnt < 0 {
		retryCnt = 0
	}

	var err error = nil
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if len(header) > 0 {
		req.Header = header[0]
	}

	var resp *http.Response
	for retryCnt >= 0 {
		if resp, err = DefaultHttpClient.Do(req); err == nil {
			return resp, nil
		}
		retryCnt--
	}
	return nil, err
}
