package gapi

import (
	"fmt"
)

// Tags represents a set of tags
type Tags map[string]bool

// NewTags creates a new Tags object from the given string slice
func NewTags(tagslice []string) Tags {
	tags := Tags{}
	tags.Set(tagslice...)
	return tags
}

// Add will add the given tags to the Tag list
func (t Tags) Add(tags ...string) {
	for _, tag := range tags {
		t[tag] = true
	}
}

// Remove will remote the given tags to the Tag list
func (t Tags) Remove(tags ...string) {
	for _, tag := range tags {
		delete(t, tag)
	}
}

// Set will delete all existing tags, and add the given tags
func (t Tags) Set(tags ...string) {
	for tag := range t {
		delete(t, tag)
	}

	t.Add(tags...)
}

// Strings will give back the tags as a string slice
func (t Tags) Strings() []string {
	tags := []string{}
	for tag := range t {
		tags = append(tags, tag)
	}

	return tags
}

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

// NewDashboard creates a new blank dashboard
func NewDashboard() *Dashboard {
	return &Dashboard{Meta: DashboardMeta{}, Model: map[string]interface{}{}}
}

// Title returns the title of the dashboard
func (d Dashboard) Title() (string, bool) {
	ititle, found := d.Model["title"]
	if !found {
		return "", false
	}

	title, ok := ititle.(string)
	return title, ok
}

// Tags returns the tags for the dashboard
func (d *Dashboard) Tags() []string {
	itagslice, found := d.Model["tags"]
	if !found {
		return []string{}
	}

	tagslice, ok := itagslice.([]string)
	if !ok {
		return []string{}
	}

	return tagslice
}

// AddTags will add the given tags to the dashboard
func (d *Dashboard) AddTags(newtags ...string) {
	tags := NewTags(d.Tags())

	for _, tag := range newtags {
		tags[tag] = true
	}

	d.Model["tags"] = tags.Strings()
}

// RemoveTags will remove the given tags to the dashboard
func (d *Dashboard) RemoveTags(deltags ...string) {
	tags := NewTags(d.Tags())

	for _, tag := range deltags {
		delete(tags, tag)
	}

	d.Model["tags"] = tags.Strings()
}

// SetTags will set the given tags on the dashboard (deleting all others)
func (d *Dashboard) SetTags(newtags ...string) {
	tags := NewTags(newtags)
	d.Model["tags"] = tags.Strings()
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

// DashboardMetas returns the dashboard metadata for the current
// organisation context.  These can then be used to get specific
// dashboards
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

// Dashboards returns the dashboards for the current org
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
