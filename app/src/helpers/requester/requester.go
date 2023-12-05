package requester

/*
	@file
	requester.go

	description:
	MakeRequest allows you to make http requests to external urls returning JSON data.
	To infer data type make sure to pass the 'target' parameter after method

	example:
	var user User
	jsonData, err := MakeRequest("example.com/user/1", "", "GET", &user)
	if err != nil {
		// log error
	}
	log(user)
*/
import (
	"bytes"
	"cloudview/app/src/api/middleware/logger"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type Requester struct {
	url     string
	data    interface{}
	method  string
	headers map[string]string
}

func MakeRequest(url string, data interface{}, method string, target interface{}, opts ...func(*Requester)) error {
	r := Requester{
		method: method,
		data:   data,
		url:    url,
	}
	// Apply functional options
	// Usually used to send optional header parameters
	for _, opt := range opts {
		opt(&r)
	}

	var body io.Reader
	switch strings.ToLower(r.method) {
	case "get":
		body = nil
	case "post":
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewBuffer(jsonData)
	default:
		return errors.New("Invalid HTTP method specified: " + method)
	}
	req, err := http.NewRequest(r.method, r.url, body)
	if err != nil {
		return err
	}

	// Set headers from Requester to the http.Request object
	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	// Perform the HTTP request
	client := &http.Client{}
	logger.Logger.Log("requester.MakeRequest: ", req.Method, req.URL.Host, req.URL.Path)
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Error("requester.MakeRequest: ERROR", err)
		return err
	}

	defer resp.Body.Close()

	/*
		To return a JSON response with type inferance if there are no errors
		in the request we need a target parameter
	*/
	jsonDataFromHttp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonDataFromHttp, target)
}

/*
Prepare the headers object to be passed to the MakeRequest function
*/
func WithHeaders(headers map[string]string) func(*Requester) {
	return func(r *Requester) {
		r.headers = headers
	}
}
