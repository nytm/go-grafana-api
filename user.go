package gapi

import (
	"bytes"
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

type UserUpdate struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Login string `json:"login,omitempty"`
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

func (c *Client) User(userId int64) (User, error) {
	var user User
	req, err := c.newRequest("GET", fmt.Sprintf("/api/users/%d", userId), nil, nil)
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

	err = json.Unmarshal(data, &user)
	return user, err
}

func (c *Client) UpdateUser(userId int64, u UserUpdate) error {
	data, err := json.Marshal(u)
	req, err := c.newRequest("PUT", fmt.Sprintf("/api/users/%d", userId), nil, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		data, _ = ioutil.ReadAll(resp.Body)
		return fmt.Errorf("status: %s body: %s", resp.Status, data)
	}
	return nil
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
