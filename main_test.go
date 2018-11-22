package gapi

import (
	"os"
	"testing"
)

var realClient *Client

func TestMain(m *testing.M) {
	var err error
	grafanaAuth := os.Getenv("GRAFANA_AUTH")
	grafanaUrl := os.Getenv("GRAFANA_URL")
	if grafanaAuth == "" {
		grafanaAuth = "admin:pwd4test"
	}
	if grafanaUrl == "" {
		grafanaUrl = "http://localhost:3000"
	}
	realClient, err = New(grafanaAuth, grafanaUrl)
	if err != nil {
		panic(err)
	}
	m.Run()
}
