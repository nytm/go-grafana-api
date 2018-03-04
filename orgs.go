package gapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	OrgUserRoleViewer = "Viewer"
	OrgUserRoleAdmin  = "Admin"
	OrgUserRoleEditor = "Editor"
)

// Org represents an Organisation object in Grafana
type Org struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

// DataSources use the given client to return the datasources
// for the organisation
func (o Org) DataSources(c *Client) ([]*DataSource, error) {
	return c.DataSourcesByOrgId(o.Id)
}

func (c *Client) AddUserToOrg(orgID int64, username, role string) error {
	data, err := json.Marshal(map[string]string{"role": role, "loginOrEmail": username})
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", fmt.Sprintf("/api/orgs/%d/user", orgID), bytes.NewReader(data))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

func (c *Client) RemoveUserFromOrg(orgID, userID int64) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/orgs/%d/users/%d", orgID, userID), nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	_, err = ioutil.ReadAll(resp.Body)
	return err
}

// Org returns the organisation with the given ID
func (c *Client) Org(id int64) (Org, error) {
	org := Org{}

	req, err := c.newRequest("GET", fmt.Sprintf("/api/orgs/%d", id), nil)
	if err != nil {
		return org, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return org, err
	}
	if resp.StatusCode != 200 {
		return org, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return org, err
	}
	err = json.Unmarshal(data, &org)
	return org, err
}

// OrgByName returns the organisation with the given name
func (c *Client) OrgByName(name string) (Org, error) {
	org := Org{}

	req, err := c.newRequest("GET", fmt.Sprintf("/api/orgs/name/%s", name), nil)
	if err != nil {
		return org, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return org, err
	}
	if resp.StatusCode != 200 {
		return org, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return org, err
	}
	err = json.Unmarshal(data, &org)
	return org, err
}

// Orgs returns all the orgs in Grafana
func (c *Client) Orgs() ([]Org, error) {
	orgs := make([]Org, 0)

	req, err := c.newRequest("GET", "/api/orgs/", nil)
	if err != nil {
		return orgs, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return orgs, err
	}
	if resp.StatusCode != 200 {
		return orgs, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return orgs, err
	}
	err = json.Unmarshal(data, &orgs)
	return orgs, err
}

// NewOrg creates an Org with the given name in Grafana
func (c *Client) NewOrg(name string) (Org, error) {
	org := Org{Name: name}
	data, err := json.Marshal(org)
	req, err := c.newRequest("POST", "/api/orgs", bytes.NewBuffer(data))
	if err != nil {
		return org, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return org, err
	}
	if resp.StatusCode != 200 {
		return org, errors.New(resp.Status)
	}

	body := struct {
		ID int64 `json:"orgId"`
	}{0}

	data, err = ioutil.ReadAll(resp.Body)
	json.Unmarshal(data, &body)
	org.Id = body.ID

	return org, err
}

// DeleteOrg deletes the given org ID from Grafana
func (c *Client) DeleteOrg(id int64) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/orgs/%d", id), nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return err
}
