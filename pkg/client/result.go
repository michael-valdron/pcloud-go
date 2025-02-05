package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type Result struct {
	Result int    `json:"result"`
	Error  string `json:"error"`
}

// ConvertToBuffer; convert http.Response.Body to bytes.Buffer
func ConvertToBuffer(resp *http.Response, err error) (*bytes.Buffer, error) {
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	defer resp.Body.Close()

	_, err = buf.ReadFrom(resp.Body)
	return buf, err
}

// CheckResult; returned error if request is failed or server returned error
func CheckResult(resp *http.Response, err error) error {
	buf, err := ConvertToBuffer(resp, err)
	if err != nil {
		return err
	}

	result := Result{}

	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		return err
	}

	if result.Result != 0 {
		return errors.New(result.Error)
	}

	return nil
}
