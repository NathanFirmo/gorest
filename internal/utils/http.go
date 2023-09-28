package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func MakeRequest(url string) ([]byte, error) {
	client := &http.Client{}

	res, err := client.Get(url)
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
