package gapi

import (
	"fmt"
)

// DashboardMeta holds dashboard metadata
type DashboardMeta struct {
	IsStarred bool     `json:"isStarred"`
	Slug      string   `json:"slug"`
	Title     string   `json:"title"`
	URI       string   `json:"uri"`
	Type      string   `json:"type"`
	Tags      []string `json:"tags"`
}

// DashboardSaveResponse represents the response from the API when
// a dashboard is saved
type DashboardSaveResponse struct {
	ID      int64  `json:"id"`
	UID     int64  `json:"uid"`
	URL     string `json:"url"`
	Status  string `json:"status"`
	Version int64  `json:"version"`
	Slug    string `json:"slug"`
}

// Dashboard represents a Grafana dashboard
type Dashboard struct {
	Meta  DashboardMeta          `json:"meta"`
	Model map[string]interface{} `json:"dashboard"`
}

// SaveDashboard saves the given dashboard model to the API
func (c *Client) SaveDashboard(model map[string]interface{}, overwrite bool) (*DashboardSaveResponse, error) {
	wrapper := map[string]interface{}{
		"dashboard": model,
		"overwrite": overwrite,
	}

	res, err := c.doJSONRequest("POST", "/api/dashboards/db", wrapper)
	if err != nil {
		return nil, err
	}

	if !res.OK() {
		return nil, res.Error()
	}

	result := &DashboardSaveResponse{}
	err = res.BindJSON(&result)
	return result, err
}

// Dashboard gets the dashboard with the given URI from Grafana
func (c *Client) Dashboard(uri string) (*Dashboard, error) {
	path := fmt.Sprintf("/api/dashboards/%s", uri)
	res, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if !res.OK() {
		return nil, res.Error()
	}

	result := &Dashboard{}
	err = res.BindJSON(&result)
	return result, err
}

func (c *Client) DashboardMetas() ([]*DashboardMeta, error) {
	res, err := c.doRequest("GET", "/api/search", nil)
	if err != nil {
		return nil, err
	}

	if !res.OK() {
		return nil, res.Error()
	}

	result := []*DashboardMeta{}
	err = res.BindJSON(&result)
	return result, err
}

func (c *Client) Dashboards() ([]*Dashboard, error) {
	dashes := []*Dashboard{}
	metas, err := c.DashboardMetas()
	if err != nil {
		return dashes, err
	}

	for _, meta := range metas {
		dash, err := c.Dashboard(meta.URI)
		if err != nil {
			return dashes, err
		}

		dashes = append(dashes, dash)
	}

	return dashes, nil
}

// DeleteDashboard will delete the dashboard with the given slug from Grafana
func (c *Client) DeleteDashboard(uri string) error {
	path := fmt.Sprintf("/api/dashboards/%s", uri)
	res, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return res.Error()
}
