package utils

import (
	"bytes"
	"io"
	"net/http"
)

var client = &http.Client{}

func SendHTTPRequest(method, url string, headers map[string]string, body []byte) (int, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, respBody, nil
}
