package json_helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func PostRequest(url string, body interface{}, result interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(respBody, result); err != nil {
		return err
	}

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
