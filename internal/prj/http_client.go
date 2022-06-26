package prj

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func DoRequest(
	client *http.Client,
	url, method string,
	payload interface{},
	headers map[string]string,
) (response *http.Response, err error) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("DoRequest json.Marshal error: %s\n", err.Error())
		return
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("DoRequest http.NewRequest error: %s\n", err.Error())
		return
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return client.Do(req)
}
