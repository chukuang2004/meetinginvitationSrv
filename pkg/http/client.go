package http

import (
	"bytes"
	"encoding/json"
	"io"
	osHttp "net/http"
	"time"
)

func Get(url string, respBody interface{}) error {

	err := custom("GET", url, nil, nil, respBody)

	return err
}

func Post(url string, reqBody, respBody interface{}) error {

	header := map[string]string{
		"Content-Type": "application/json",
	}
	err := custom("POST", url, header, reqBody, respBody)

	return err
}

func CustomPost(url string, reqHeader map[string]string, reqBody, respBody interface{}) error {

	err := custom("POST", url, reqHeader, reqBody, respBody)

	return err
}

func custom(method, url string, reqHeader map[string]string, reqBody, respBody interface{}) error {

	client := &osHttp.Client{
		Timeout: time.Second * 5,
	}

	tmp, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := osHttp.NewRequest(method, url, bytes.NewBuffer(tmp))
	if err != nil {
		return err
	}

	for k, v := range reqHeader {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, respBody)

	return err
}
