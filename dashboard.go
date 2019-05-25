package gapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

type DashboardMeta struct {
	IsStarred   bool   `json:"isStarred"`
	Slug        string `json:"slug"`
	Folder      int64  `json:"folderId"`
	FolderTitle string `json:"folderTitle"`
}

// DashboardSaveResponse grafana response for create dashboard
type DashboardSaveResponse struct {
	Slug    string `json:"slug"`
	ID      int64  `json:"id"`
	UID     string `json:"uid"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Version int64  `json:"version"`
}

type Dashboard struct {
	Meta      DashboardMeta          `json:"meta"`
	Model     map[string]interface{} `json:"dashboard"`
	Folder    int64                  `json:"folderId"`
	Overwrite bool                   `json:overwrite`
}

// Dashboards represent json returned by search API
type Dashboards struct {
	ID          int64  `json:"id"`
	UID         string `json:"uid"`
	Title       string `json:"title"`
	URI         string `json:"uri"`
	URL         string `json:"url"`
	Starred     bool   `json:"isStarred"`
	FolderID    int64  `json:"folderId"`
	FolderUID   string `json:"folderUid"`
	FolderTitle string `json:"folderTitle"`
}

// DashboardDeleteResponse grafana response for delete dashboard
type DashboardDeleteResponse struct {
	Title string `json:title`
}

// Deprecated: use NewDashboard instead
func (c *Client) SaveDashboard(model map[string]interface{}, overwrite bool) (*DashboardSaveResponse, error) {
	wrapper := map[string]interface{}{
		"dashboard": model,
		"overwrite": overwrite,
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest("POST", "/api/dashboards/db", nil, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		data, _ = ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode, data)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &DashboardSaveResponse{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (c *Client) NewDashboard(dashboard Dashboard) (*DashboardSaveResponse, error) {
	data, err := json.Marshal(dashboard)
	if err != nil {
		return nil, err
	}
	req, err := c.newRequest("POST", "/api/dashboards/db", nil, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &DashboardSaveResponse{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// SearchDashboard search a dashboard in Grafana
func (c *Client) SearchDashboard(query string, folderID string) ([]Dashboards, error) {
	dashboards := make([]Dashboards, 0)
	path := "/api/search"

	params := url.Values{}
	params.Add("type", "dash-db")
	params.Add("query", query)
	params.Add("folderIds", folderID)

	req, err := c.newRequest("GET", path, params, nil)
	if err != nil {
		return dashboards, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return dashboards, err
	}
	if resp.StatusCode != 200 {
		return dashboards, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return dashboards, err
	}

	err = json.Unmarshal(data, &dashboards)

	return dashboards, err
}

// GetDashboard get a dashboard by UID
func (c *Client) GetDashboard(uid string) (*Dashboard, error) {
	path := fmt.Sprintf("/api/dashboards/uid/%s", uid)
	req, err := c.newRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Dashboard{}
	err = json.Unmarshal(data, &result)
	result.Folder = result.Meta.Folder
	if os.Getenv("GF_LOG") != "" {
		log.Printf("got back dashboard response  %s", data)
	}
	return result, err
}

// Deprecated: use GetDashboard instead
func (c *Client) Dashboard(slug string) (*Dashboard, error) {
	path := fmt.Sprintf("/api/dashboards/db/%s", slug)
	req, err := c.newRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Dashboard{}
	err = json.Unmarshal(data, &result)
	result.Folder = result.Meta.Folder
	if os.Getenv("GF_LOG") != "" {
		log.Printf("got back dashboard response  %s", data)
	}
	return result, err
}

// DeleteDashboard deletes a grafana dashoboard
func (c *Client) DeleteDashboard(uid string) (string, error) {
	deleted := &DashboardDeleteResponse{}
	path := fmt.Sprintf("/api/dashboards/uid/%s", uid)
	req, err := c.newRequest("DELETE", path, nil, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(data, &deleted)
	if err != nil {
		return "", err
	}
	return deleted.Title, nil
}
