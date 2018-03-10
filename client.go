package gapi

import (
	"encoding/base64"
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

// Client represents a Grafana API client
type Client struct {
	bearerAuth string
	basicAuth  string
	baseURL    url.URL
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

func (c *Client) parseAuth(auth string) {
	if strings.Contains(auth, ":") {
		c.basicAuth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
		split := strings.Split(auth, ":")
		c.baseURL.User = url.UserPassword(split[0], split[1])
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

	if c.bearerAuth != "" {
		req.Header.Add("Authorization", c.bearerAuth)
	}

	if os.Getenv("GF_LOG") != "" {
		if body == nil {
			log.Println("request to ", url.String(), "with no body data")
		} else {
			data, _ := ioutil.ReadAll(body)
			log.Println("request to ", url.String(), "with body data", string(data))
		}
	}

	req.Header.Add("Content-Type", "application/json")
	return req, err
}
