package json_helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
)

func PostRequest(url string, body interface{}, result interface{}) error {

	log.Printf("eth-blocks-requester::PostRequest::url: %v\n", url)
	log.Printf("eth-blocks-requester::PostRequest::body: %v\n", body)

	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Printf("eth-blocks-requester::PostRequest::ERROR::Failed to marshal JSON: %v\n", err)
		return err
	}

	log.Printf("eth-blocks-requester::PostRequest::jsonData: %v\n", jsonData)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("eth-blocks-requester::ERROR::Failed to create request: %v\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("eth-blocks-requester::ERROR::Failed to send request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("eth-blocks-requester::ERROR::Failed to read response body: %v\n", err)
		return err
	}

	log.Printf("eth-blocks-requester::PostRequest::respBody: %v\n", respBody)
	log.Printf("eth-blocks-requester::PostRequest::[]byte(respBody): %v\n", respBody)
	log.Printf("eth-blocks-requester::PostRequest::string(respBody): %v\n", string(respBody))

	if err := json.Unmarshal(respBody, result); err != nil {
		log.Printf("eth-blocks-requester::ERROR::Failed to unmarshal response: %v\n", err)
		return err
	}

	log.Printf("eth-blocks-requester::PostRequest::result: %v\n", result)

	if hasError(result) {
		return fmt.Errorf("eth-blocks-requester::PostRequest::Response contains an error field")
	}

	return nil
}

func hasError(inputStruct interface{}) bool {
	value := reflect.ValueOf(inputStruct)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return false
	}
	errorField := value.FieldByName("Error")
	return errorField.IsValid() && !errorField.IsZero()
}
