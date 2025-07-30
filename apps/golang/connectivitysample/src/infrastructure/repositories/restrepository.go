package repositories

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Send a POST Request
func SendPostRequest(url string, body string, contentType string) error {

	response, err := http.Post(url,
		contentType,
		bytes.NewBuffer([]byte(body)))

	if err != nil {
		log.Printf("Error making a POST request: %v", err)
		return err
	}
	if response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		errorMessage := fmt.Sprintf("Error when doing the POST request. Response status:'%s' Body:'%s'", response.Status, string(body))
		log.Printf(errorMessage)
		return errors.New(errorMessage)
	}
	return nil
}
