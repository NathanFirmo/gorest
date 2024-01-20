package utils

import (
	"encoding/json"
	// "fmt"
	"io"
	"net/http"
	"strings"
)

func MakeRequest(method string, url string, rawHeaders string, rawBody string) ([]byte, error) {
	req, _ := http.NewRequest(method, url, strings.NewReader(rawBody))
	headers := strings.Split(strings.TrimSpace(rawHeaders), "\n")

	for _, h := range headers {
		parsedHeader := strings.Split(h, ":")
		if len(parsedHeader) == 1 {
			continue
		}
		req.Header.Add(parsedHeader[0], parsedHeader[1])
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		return nil, err
	}

	indentedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return nil, err
	}

	return indentedJSON, nil
}
