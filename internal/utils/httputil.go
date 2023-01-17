package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gookit/slog"
)

// encodeResponseWithStatusWithError ...
func encodeResponseWithStatusWithError(w http.ResponseWriter, val interface{}, statusCode int) error {
	if val == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(val)
}

// EncodeResponseWithStatus ...
func EncodeResponseWithStatus(w http.ResponseWriter, val interface{}, statusCode int) {
	err := encodeResponseWithStatusWithError(w, val, statusCode)

	// log the error
	if err != nil {
		slog.Errorf("Error writing response, error: %w, type: %T", err, val)
	}
}

// DecodeRequest decode request into the struct.
func DecodeRequest(r *http.Request, val interface{}) error {
	if r.Body == nil {
		slog.Errorf("Request body could not be read, error: %q", "EMPTY_REQUEST_BODY")
		return fmt.Errorf("EMPTY_REQUEST_BODY")
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Errorf("%s, error: %q", err.Error(), "EMPTY_REQUEST_BODY")
		return fmt.Errorf("REQUEST_READING_FAILED")
	}

	err = json.Unmarshal(data, val)
	if err != nil {
		slog.Errorf("%s, error: %q, data: %s", err.Error(), "REQUEST_PARSING_FAILED", string(data))
		return fmt.Errorf("REQUEST_PARSING_FAILED")
	}
	return nil
}

// CreateRequest creates a http request.
func CreateRequest(method, url, token string, data []byte) (req *http.Request, err error) {
	// new http request
	slog.Debugf("%s -> %v : %v", method, url, string(data))
	req, err = http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return req, fmt.Errorf("failed to create http request [%s]%s: %w", method, url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	return req, nil
}

// SendRequest sends a http request and validates if the response was success, pass 0 for successStatusCode if validation is not required
func SendRequest(req *http.Request, successStatusCode int) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	var client http.Client

	// make request
	resp, err = client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("failed to make http request [%s] %s: %w", req.Method, req.URL.Path, err)
	}

	// check status code  if passed
	if successStatusCode > 0 && resp.StatusCode != successStatusCode {
		logUnexpectedResponseOutput(*resp)
		return resp, fmt.Errorf("received failure status code '%v'", resp.StatusCode)
	}

	return resp, err
}

func logUnexpectedResponseOutput(resp http.Response) error {
	slog.Errorf("Unexpected error code %v", resp.StatusCode)
	tmp, _ := io.ReadAll(resp.Body)
	var responseBody interface{}
	err := json.Unmarshal(tmp, &responseBody)
	if err != nil {
		slog.Errorf("Failed to unmarshal response %v", err)
		return err
	} else {
		slog.Errorf("Unexpected response body: %v", responseBody)
	}
	return nil
}
