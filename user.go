package gapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

// User represents a Grafana user
type User struct {
	Id      int64
	Email   string
	Name    string
	Login   string
	OrgId   string
	IsAdmin bool
}

// Users returns all the users from Grafana
func (c *Client) Users() ([]User, error) {
	users := make([]User, 0)
	req, err := c.newRequest("GET", "/api/users", nil)
	if err != nil {
		return users, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return users, err
	}
	if resp.StatusCode != 200 {
		return users, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return users, err
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		return users, err
	}
	return users, err
}

func (c *Client) UpdateDefaultOrg(userID, orgID int64) error {
	req, err := c.newRequest("POST", fmt.Sprintf("/api/users/%d/using/%d", userID, orgID), nil)
	if err != nil {
		return err
	}

	req.URL.User = nil
	log.Printf("%+v", req.URL)
	req.Header.Set("Authorization", "Basic "+c.authBasic)
	log.Printf("%+v", req.Header)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}
