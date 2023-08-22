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
	jsonData, err := json.Marshal(body)
	if err != nil {
		log.Printf("eth-transactions-requester::ERROR::Failed to marshal JSON: %v\n", err)
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("eth-transactions-requester::ERROR::Failed to create request: %v\n", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("eth-transactions-requester::ERROR::Failed to send request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("eth-transactions-requester::ERROR::Failed to read response body: %v\n", err)
		return err
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		log.Printf("eth-transactions-requester::ERROR::Failed to unmarshal response: %v\n", err)
		return err
	}

	if hasError(result) {
		log.Printf("eth-transactions-requester::PostRequest::ERROR:: %v\n", result)
		return fmt.Errorf("eth-transactions-requester::PostRequest::ERROR:: %v/n", "Response contains an error field")
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
