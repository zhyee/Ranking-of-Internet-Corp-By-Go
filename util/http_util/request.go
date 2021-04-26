package http_util

import "net/http"

/**
带重试的Get请求，最多尝试 1 + retryCnt 次请求，最少尝试1次请求
 */
func GetWithRetry(url string, retryCnt int) (*http.Response, error) {

	if retryCnt < 0 {
		retryCnt = 0
	}

	var resp *http.Response = nil
	var err error = nil
	for i := 0; i <= retryCnt; i++ {
		resp, err = http.Get(url)
		if err == nil && resp != nil {
			return resp, nil
		}
	}
	return resp, err
}
