package gapi

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var ErrNotFound = errors.New("not found")
var ErrNotImplemented = errors.New("not implemented")

// Client represents a Grafana API client
type Client struct {
	bearerAuth     string
	basicAuth      string
	baseURL        url.URL
	LastStatusCode int
	*http.Client
}

// New creates a new grafana client
// auth can be in user:pass format, or it can be an api key
func New(auth, baseURL string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		Client:  &http.Client{},
		baseURL: *u,
	}

	c.parseAuth(auth)

	return c, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	res, err := c.Client.Do(req)
	if err != nil {
		return res, err
	}

	logResponse(res)

	c.LastStatusCode = res.StatusCode
	return res, err
}

func (c *Client) parseAuth(auth string) {
	if strings.Contains(auth, ":") {
		c.basicAuth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
	} else {
		c.bearerAuth = fmt.Sprintf("Bearer %s", auth)
	}
}

func (c *Client) newRequest(method, requestPath string, body io.Reader) (*http.Request, error) {
	url := c.baseURL.String() + requestPath
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return req, err
	}

	if c.basicAuth != "" {
		req.Header.Add("Authorization", c.basicAuth)
	}

	if c.bearerAuth != "" {
		req.Header.Add("Authorization", c.bearerAuth)
	}

	req.Header.Add("Content-Type", "application/json")

	logRequest(req)

	return req, err
}

func logRequest(req *http.Request) {
	if os.Getenv("GF_LOG") == "" {
		return
	}

	fmt.Println("\nHTTP/1.1", req.Method, req.URL)
	req.Header.Write(os.Stdout)

	if req.Body != nil {
		data, _ := ioutil.ReadAll(req.Body)
		fmt.Println(string(data))
	}

	fmt.Println("")
}

func logResponse(res *http.Response) {
	if os.Getenv("GF_LOG") == "" {
		return
	}

	fmt.Println("\nRESPONSE HEADERS:")
	res.Header.Write(os.Stdout)

	if os.Getenv("GF_LOG") != "2" {
		return
	}

	buf1 := bytes.NewBuffer([]byte{})
	buf2 := bytes.NewBuffer([]byte{})
	mw := io.MultiWriter(buf1, buf2)
	_, _ = io.Copy(mw, res.Body)
	res.Body = ioutil.NopCloser(bytes.NewReader(buf1.Bytes()))
	fmt.Println("")
	fmt.Println(string(buf2.Bytes()))
	fmt.Println("")
}
