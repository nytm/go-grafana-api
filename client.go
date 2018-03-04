package gapi

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// Client represents a Grafana API client
type Client struct {
	key       string
	baseURL   url.URL
	authBasic string
	*http.Client
}

// New creates a new grafana client
// auth can be in user:pass format, or it can be an api key
func New(auth, baseURL string) (*Client, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	key := ""
	authBasic := ""
	if strings.Contains(auth, ":") {
		authBasic = base64.StdEncoding.EncodeToString([]byte(auth))
		split := strings.Split(auth, ":")
		u.User = url.UserPassword(split[0], split[1])
	} else {
		key = fmt.Sprintf("Bearer %s", auth)
	}

	return &Client{
		key,
		*u,
		authBasic,
		&http.Client{},
	}, nil
}

func (c *Client) newRequest(method, requestPath string, body io.Reader) (*http.Request, error) {
	url := c.baseURL
	url.Path = path.Join(url.Path, requestPath)
	req, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return req, err
	}
	if c.key != "" {
		req.Header.Add("Authorization", c.key)
	}

	if os.Getenv("GF_LOG") != "" {
		if body == nil {
			log.Println("request to ", url.String(), "with no body data")
		} else {
			log.Println("request to ", url.String(), "with body data", body.(*bytes.Buffer).String())
		}
	}

	req.Header.Add("Content-Type", "application/json")
	return req, err
}
