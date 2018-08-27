package gapi

import (
	"fmt"
)

// AdminCreateUserForm is used to create a new user
type AdminCreateUserForm struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"Required"`
}

// CreateUserForm will create a user from the given form
func (c *Client) CreateUserForm(settings AdminCreateUserForm) error {
	res, err := c.doJSONRequest("POST", "/api/admin/users", settings)
	if err != nil {
		return err
	}

	return res.Error()
}

// DeleteUser deletes the user with the given ID from Grafana
func (c *Client) DeleteUser(id int64) error {
	res, err := c.doRequest("DELETE", fmt.Sprintf("/api/admin/users/%d", id), nil)
	if err != nil {
		return err
	}

	return res.Error()
}

// Stats will get the stats from the API
func (c *Client) Stats() (map[string]int64, error) {
	v := map[string]int64{}
	res, err := c.doRequest("GET", "/api/admin/stats", nil)
	if err != nil {
		return v, err
	}

	if !res.OK() {
		return v, res.Error()
	}

	err = res.BindJSON(&v)
	return v, err
}

// FrontEndSettings will get the front end settings from the API
func (c *Client) FrontEndSettings() (map[string]interface{}, error) {
	v := map[string]interface{}{}
	res, err := c.doRequest("GET", "/api/frontend/settings", nil)
	if err != nil {
		return v, err
	}

	if !res.OK() {
		return v, res.Error()
	}

	err = res.BindJSON(&v)
	return v, err
}
