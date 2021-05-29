package httclient

import (
	"bytes"
	"net/http"
)

var CallRestAPI = func(method, url string, headers map[string]string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for headerKey, headerVal := range headers {
		req.Header.Set(headerKey, headerVal)
	}

	var resp *http.Response

	client := &http.Client{}

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
