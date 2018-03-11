package gapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/grafana/grafana/pkg/api/dtos"
)

// CreateUserForm will create a user from the given form
func (c *Client) CreateUserForm(settings dtos.AdminCreateUserForm) error {
	data, err := json.Marshal(settings)
	req, err := c.newRequest("POST", "/api/admin/users", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return err
}

// DeleteUser deletes the user with the given ID from Grafana
func (c *Client) DeleteUser(id int64) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/admin/users/%d", id), nil)
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

// Stats will get the stats from the API
func (c *Client) Stats() (map[string]int64, error) {
	v := map[string]int64{}
	req, err := c.newRequest("GET", "/api/admin/stats", nil)
	if err != nil {
		return v, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return v, err
	}

	if resp.StatusCode != 200 {
		return v, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(data, &v)
	return v, err
}

// FrontEndSettings will get the front end settings from the API
func (c *Client) FrontEndSettings() (map[string]interface{}, error) {
	v := map[string]interface{}{}
	req, err := c.newRequest("GET", "/api/frontend/settings", nil)
	if err != nil {
		return v, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return v, err
	}

	if resp.StatusCode != 200 {
		return v, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return v, err
	}

	err = json.Unmarshal(data, &v)
	return v, err
}
