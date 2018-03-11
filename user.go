package gapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"github.com/grafana/grafana/pkg/api/dtos"
)

// User represents a Grafana user
type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	OrgId    string `json:"org_id"`
	IsAdmin  bool   `json:"isGrafanaAdmin"` // TODO: handle isAdmin returned from /api/users
	Password string `json:"password,omitempty"`
}

type Users []User

func (users Users) FindByEmail(email string) (User, bool) {
	for _, u := range users {
		if u.Email == email {
			return u, true
		}
	}

	return User{}, false
}

func (users Users) FindIndexByEmail(email string) (int, bool) {
	for i, u := range users {
		if u.Email == email {
			return i, true
		}
	}

	return 0, false
}

func (u User) Using(c *Client, orgID int64) error {
	req, err := c.newRequest("POST", fmt.Sprintf("/api/users/%d/using/%d", u.Id, orgID), nil)
	if err != nil {
		return err
	}

	req.URL.User = nil
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
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

func (c *Client) User(id int64) (User, error) {
	user := User{}
	req, err := c.newRequest("GET", fmt.Sprintf("/api/users/%d", id), nil)
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
	if err != nil {
		return user, err
	}
	user.Id = id
	return user, err
}

// NewUser creates a new user by wrapping the CreateUserForm method to
// avoiding requiring a dependency on Grafana code in your code
func (c *Client) NewUser(u User) error {
	form := dtos.AdminCreateUserForm{}
	form.Password = u.Password
	form.Email = u.Email
	form.Name = u.Name
	form.Login = u.Login

	return c.CreateUserForm(form)
}

// SwitchUserOrg will switch the current organisation (uses basic auth)
func (c *Client) SwitchUserOrg(userID, orgID int64) error {
	req, err := c.newRequest("POST", fmt.Sprintf("/api/users/%d/using/%d", userID, orgID), nil)
	if err != nil {
		return err
	}

	req.URL.User = nil
	resp, err := c.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	log.Println("hihiihi", string(data))
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}

// SwitchCurrentUserOrg will switch the current organisation of the signed in user
func (c *Client) SwitchCurrentUserOrg(orgID int64) error {
	req, err := c.newRequest("POST", fmt.Sprintf("/api/user/using/%d", orgID), nil)
	if err != nil {
		return err
	}

	req.URL.User = nil
	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	return nil
}

// UserByEmail will find a user by their email address
func (c *Client) UserByEmail(email string) (User, error) {
	user := User{}

	values := url.Values{}
	values.Set("loginOrEmail", email)

	req, err := c.newRequest("GET", "/api/users/lookup?"+values.Encode(), nil)
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
	if err != nil {
		return user, err
	}
	return user, err
}
