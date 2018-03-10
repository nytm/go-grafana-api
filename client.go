package gapi

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
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
	if os.Getenv("GF_LOG") == "2" {
		log.Println("===> GAPI: request headers:")
		res.Request.Header.Write(os.Stdout)
		log.Println("===> GAPI: response headers:")
		res.Header.Write(os.Stdout)

		buf1 := bytes.NewBuffer([]byte{})
		buf2 := bytes.NewBuffer([]byte{})
		mw := io.MultiWriter(buf1, buf2)
		_, _ = io.Copy(mw, res.Body)
		res.Body = ioutil.NopCloser(bytes.NewReader(buf1.Bytes()))
		log.Println("===> GAPI: response body:", string(buf2.Bytes()))
	}

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
	url := c.baseURL
	url.Path = path.Join(url.Path, requestPath)
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return req, err
	}

	if c.basicAuth != "" {
		req.Header.Add("Authorization", c.basicAuth)
	}

	if c.bearerAuth != "" {
		req.Header.Add("Authorization", c.bearerAuth)
	}

	if os.Getenv("GF_LOG") != "" {
		log.Println("===> GAPI: request to ", url.String(), "with no body data")
		if body != nil {
			data, _ := ioutil.ReadAll(body)
			log.Println("===> GAPI: request body:", string(data))
		}
	}

	req.Header.Add("Content-Type", "application/json")

	return req, err
}
