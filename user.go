package gapi

import (
	"fmt"
	"net/url"
)

// User represents a Grafana user
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Login    string `json:"login"`
	OrgID    string `json:"org_id"`
	IsAdmin  bool   `json:"isGrafanaAdmin"` // TODO: handle isAdmin returned from /api/users
	Password string `json:"password,omitempty"`
}

// Users is a collection of user models
type Users []*User

// FindByEmail returns the user with the given email from a
// collection of users, and a false if it was not found
func (users Users) FindByEmail(email string) (*User, bool) {
	for _, u := range users {
		if u.Email == email {
			return u, true
		}
	}

	return &User{}, false
}

// FindIndexByEmail is like FindByEmail but it returns the index
func (users Users) FindIndexByEmail(email string) (int, bool) {
	for i, u := range users {
		if u.Email == email {
			return i, true
		}
	}

	return 0, false
}

// FindByLogin returns the user with the given email from a
// collection of users, and a false if it was not found
func (users Users) FindByLogin(login string) (*User, bool) {
	for _, u := range users {
		if u.Login == login {
			return u, true
		}
	}

	return &User{}, false
}

// FindIndexByLogin is like FindByEmail but it returns the index
func (users Users) FindIndexByLogin(login string) (int, bool) {
	for i, u := range users {
		if u.Login == login {
			return i, true
		}
	}

	return 0, false
}

// SwitchOrg will change the current org context for the user
func (u User) SwitchOrg(c *Client, orgID int64) error {
	return c.SwitchUserOrg(u.ID, orgID)
}

// MakeGlobalAdmin assigns the user to all orgs with an Admin role
func (u User) MakeGlobalAdmin(c *Client) error {
	return u.AddToAllOrgs(c, OrgUserRoleAdmin)
}

// MakeGlobalEditor assigns the user to all orgs with a Editor role
func (u User) MakeGlobalEditor(c *Client) error {
	return u.AddToAllOrgs(c, OrgUserRoleEditor)
}

// MakeGlobalViewer assigns the user to all orgs with a Viewer role
func (u User) MakeGlobalViewer(c *Client) error {
	return u.AddToAllOrgs(c, OrgUserRoleViewer)
}

// RemoveFromAllOrgs will remove the user from all the orgs that they
// have a current role in
func (u User) RemoveFromAllOrgs(c *Client) error {
	orgs, err := c.Orgs()
	if err != nil {
		return err
	}

	for _, org := range orgs {
		ousers, err := org.Users(c)
		if err != nil {
			return err
		}

		u, ok := OrgUsers(ousers).FindByLogin(u.Login)
		if !ok {
			continue
		}

		if err := org.RemoveUser(c, u.ID); err != nil {
			return err
		}
	}

	return nil
}

// AddToAllOrgs will add the user to all orgs with the given role
func (u User) AddToAllOrgs(c *Client, role string) error {
	orgs, err := c.Orgs()
	if err != nil {
		return err
	}

	for _, org := range orgs {
		err := org.AddUser(c, u.Login, role)

		if err != nil && err != ErrConflict {
			return err
		}

		if err != nil && err == ErrConflict {
			ousers, err := org.Users(c)
			if err != nil {
				return err
			}

			ouser, ok := OrgUsers(ousers).FindByLogin(u.Login)
			if !ok {
				return fmt.Errorf("Conflict occured while assigning %s to %s, but user is not found in that org", u.Login, org.Name)
			}

			if err := ouser.UpdateRole(c, role); err != nil {
				return fmt.Errorf("unable to update role for %s on %s: %s", u.Login, org.Name, err)
			}
		}
	}

	return nil
}

// Users returns all the users from Grafana
func (c *Client) Users() ([]*User, error) {
	users := make([]*User, 0)
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

// User returns the user with the given id
func (c *Client) User(id int64) (*User, error) {
	user := &User{}
	res, err := c.doRequest("GET", fmt.Sprintf("/api/users/%d", id), nil)
	if err != nil {
		return user, err
	}

	if !res.OK() {
		return user, res.Error()
	}

	err = res.BindJSON(&user)
	user.ID = id
	return user, err
}

// NewUser is DEPRECATED
func (c *Client) NewUser(u User) error {
	return c.CreateUser(u)
}

// CreateUser creates a new user by wrapping the CreateUserForm method to
// avoiding requiring a dependency on Grafana code in your code
func (c *Client) CreateUser(u User) error {
	form := AdminCreateUserForm{}
	form.Password = u.Password
	form.Email = u.Email
	form.Name = u.Name
	form.Login = u.Login

	return c.CreateUserForm(form)
}

// SaveUser will save the given user to the API
func (c *Client) SaveUser(u *User) error {
	res, err := c.doRequest("PUT", fmt.Sprintf("/api/users/%d", u.ID), nil)
	if err != nil {
		return err
	}

	return res.Error()
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

// SetUserAdmin will set the given user ID as an admin
func (c *Client) SetUserAdmin(id int64, admin bool) error {
	body := map[string]bool{"isGrafanaAdmin": admin}
	res, err := c.doJSONRequest("PUT", fmt.Sprintf("/api/admin/users/%d/permissions", id), body)
	if err != nil {
		return err
	}

	return res.Error()
}

// UserByEmail will find a user by their email address
func (c *Client) UserByEmail(email string) (*User, error) {
	user := &User{}

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

// UserByLogin will find a user by their login
func (c *Client) UserByLogin(login string) (*User, error) {
	return c.UserByEmail(login)
}
