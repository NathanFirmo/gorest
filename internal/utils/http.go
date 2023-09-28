package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func MakeRequest(url string) ([]byte, error) {
	// Create a new HTTP client.
	client := &http.Client{}

	// Make a GET request to the Google homepage.
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	// Close the response body when we're done with it.
	defer res.Body.Close()

	// Read the response body into a byte array.
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Create a map or struct to unmarshal the JSON data
	var jsonData interface{}
	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		return nil, err
	}

	// Marshal the JSON data with indentation
	indentedJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return nil, err
	}

	return indentedJSON, nil
}
