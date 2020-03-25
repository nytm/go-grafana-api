package gapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
)

type User struct {
	Id       int64  `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	IsAdmin  bool   `json:"isAdmin,omitempty"`
}

func (c *Client) Users() ([]User, error) {
	users := make([]User, 0)
	req, err := c.newRequest("GET", "/api/users", nil, nil)
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

func (c *Client) UserByEmail(email string) (User, error) {
	user := User{}
	query := url.Values{}
	query.Add("loginOrEmail", email)
	req, err := c.newRequest("GET", "/api/users/lookup", query, nil)
	if err != nil {
		return user, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return user, err
	}
	if resp.StatusCode != 200 {
		return user, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}
	tmp := struct {
		Id       int64  `json:"id,omitempty"`
		Email    string `json:"email,omitempty"`
		Name     string `json:"name,omitempty"`
		Login    string `json:"login,omitempty"`
		Password string `json:"password,omitempty"`
		IsAdmin  bool   `json:"isGrafanaAdmin,omitempty"`
	}{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return user, err
	}
	user = User(tmp)
	return user, err
}

// CurrentUser returns user info about the currently-logged-in user
// https://grafana.com/docs/grafana/latest/http_api/user/#actual-user
func (c *Client) CurrentUser() (User, error) {
	user := User{}
	req, err := c.newRequest("GET", "/api/user", nil, nil)
	if err != nil {
		return user, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return user, err
	}
	if resp.StatusCode != 200 {
		return user, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}
	tmp := struct {
		Id       int64  `json:"id,omitempty"`
		Email    string `json:"email,omitempty"`
		Name     string `json:"name,omitempty"`
		Login    string `json:"login,omitempty"`
		Password string `json:"password,omitempty"`
		IsAdmin  bool   `json:"isGrafanaAdmin,omitempty"`
	}{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return user, err
	}
	user = User(tmp)
	return user, err
}

type OrgMembership struct {
	Id   int64  `json:"orgId"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// CurrentUserOrgs returns org membership for the currently-logged-in user
// https://grafana.com/docs/grafana/latest/http_api/user/#organizations-of-the-actual-user
func (c *Client) CurrentUserOrgs() ([]OrgMembership, error) {
	orgs := make([]OrgMembership, 0)

	req, err := c.newRequest("GET", "/api/user/orgs/", nil, nil)
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

// SwitchUserOrganization changes the user's currently-selected organization
// https://grafana.com/docs/grafana/latest/http_api/user/#switch-user-context-for-a-specified-user
func (c *Client) SwitchUserOrganization(userID, orgID int64) error {
	req, err := c.newRequest("POST", fmt.Sprintf("/api/users/%d/using/%d", userID, orgID), nil, nil)
	if err != nil {
		return err
	}

	response, err := c.Do(req)
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return errors.New(response.Status)
	}
	return nil
}
