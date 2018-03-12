package gapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewResponse returns a new grafana API response
func NewResponse(res *http.Response, rerr error) *Response {
	var data []byte

	if res.Body != nil {
		data, _ = ioutil.ReadAll(res.Body)
	}

	return &Response{
		res,
		data,
		rerr,
	}
}

// Response is an API response
type Response struct {
	*http.Response
	data []byte
	err  error
}

// OK is true if there is no error
func (res *Response) OK() bool {
	return res.Error() == nil
}

// BindJSON unmarshals the body into the given interface
func (res *Response) BindJSON(v interface{}) error {
	return json.Unmarshal(res.data, v)
}

// Message returns the message from the
func (res *Response) Message() string {
	data := struct {
		Msg string `json:"message"`
	}{}
	res.BindJSON(&data)
	return data.Msg
}

func (res *Response) Error() error {
	if res.err != nil {
		return res.err
	}

	switch res.StatusCode {
	case 200:
		return nil
	case 404:
		return ErrNotFound
	case 409:
		return ErrConflict
	default:
		return fmt.Errorf(res.Status)
	}
}
