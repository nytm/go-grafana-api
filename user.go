package gapi

import (
	"fmt"
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

func (u User) SwitchOrg(c *Client, orgID int64) error {
	return c.SwitchUserOrg(u.Id, orgID)
}

// Users returns all the users from Grafana
func (c *Client) Users() ([]User, error) {
	users := make([]User, 0)
	res, err := c.doRequest("GET", "/api/users", nil)
	if err != nil {
		return users, err
	}

	if !res.OK() {
		return users, res.Error()
	}

	err = res.BindJSON(&users)
	return users, err
}

func (c *Client) User(id int64) (User, error) {
	user := User{}
	res, err := c.doRequest("GET", fmt.Sprintf("/api/users/%d", id), nil)
	if err != nil {
		return user, err
	}

	if !res.OK() {
		return user, res.Error()
	}

	err = res.BindJSON(&user)
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
	res, err := c.doRequest("POST", fmt.Sprintf("/api/users/%d/using/%d", userID, orgID), nil)
	if err != nil {
		return err
	}

	return res.Error()
}

// SwitchCurrentUserOrg will switch the current organisation of the signed in user
func (c *Client) SwitchCurrentUserOrg(orgID int64) error {
	res, err := c.doRequest("POST", fmt.Sprintf("/api/user/using/%d", orgID), nil)
	if err != nil {
		return err
	}

	return res.Error()
}

// UserByEmail will find a user by their email address
func (c *Client) UserByEmail(email string) (User, error) {
	user := User{}

	values := url.Values{}
	values.Set("loginOrEmail", email)

	res, err := c.doRequest("GET", "/api/users/lookup?"+values.Encode(), nil)
	if err != nil {
		return user, err
	}

	if !res.OK() {
		return user, res.Error()
	}

	err = res.BindJSON(&user)
	return user, err
}
