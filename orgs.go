package gapi

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	// OrgUserRoleViewer is the readonly role
	OrgUserRoleViewer = "Viewer"
	// OrgUserRoleAdmin is the admin role
	OrgUserRoleAdmin = "Admin"
	// OrgUserRoleEditor is the editing role
	OrgUserRoleEditor = "Editor"
)

// OrgUser is a user of the org
type OrgUser struct {
	User
	Role  string `json:"role"`
	OrgID int64  `json:"org_id"`
}

// OrgUsers is a collection of Org user models
type OrgUsers []*OrgUser

// Users returns the user objects from a collection of org users
func (ousers OrgUsers) Users() []*User {
	users := []*User{}
	for _, ou := range ousers {
		users = append(users, &ou.User)
	}
	return users
}

// Org represents an Organisation object in Grafana
type Org struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (o Org) String() string {
	return o.Name
}

// DataSources use the given client to return the datasources
// for the organisation
func (o Org) DataSources(c *Client) ([]*DataSource, error) {
	return c.DataSourcesByOrgID(o.ID)
}

// AddUser will add a user to the organisation
func (o Org) AddUser(c *Client, username, role string) error {
	role = AutoFixRole(role)

	if !UserRoleIsValid(role) {
		return fmt.Errorf("invalid role name: %s", role)
	}

	acl := map[string]string{"role": role, "loginOrEmail": username}

	res, err := c.doJSONRequest("POST", fmt.Sprintf("/api/orgs/%d/users", o.ID), acl)
	if err != nil {
		return err
	}

	return res.Error()
}

// Users use the given client to return the users
// for the organisation
func (o Org) Users(c *Client) ([]*OrgUser, error) {
	ousers := []*OrgUser{}

	res, err := c.doRequest("GET", fmt.Sprintf("/api/orgs/%d/users", o.ID), nil)
	if err != nil {
		return ousers, err
	}

	if !res.OK() {
		return ousers, res.Error()
	}

	err = res.BindJSON(&ousers)
	return ousers, err
}

// RemoveUser removes the user from the organisation
func (o Org) RemoveUser(c *Client, userID int64) error {
	res, err := c.doRequest("DELETE", fmt.Sprintf("/api/orgs/%d/users/%d", o.ID, userID), nil)
	if err != nil {
		return err
	}

	return res.Error()
}

// Org returns the organisation with the given ID
func (c *Client) Org(id int64) (Org, error) {
	org := Org{}

	res, err := c.doRequest("GET", fmt.Sprintf("/api/orgs/%d", id), nil)
	if err != nil {
		return org, err
	}

	if !res.OK() {
		return org, res.Error()
	}

	err = res.BindJSON(&org)
	return org, err
}

// OrgByName returns the organisation with the given name
func (c *Client) OrgByName(name string) (Org, error) {
	org := Org{}

	// the normal query escape replaces spaces with the plus symbol
	// grafana API does not like that, use %20 instead as per API docs
	name = url.QueryEscape(name)
	name = strings.Replace(name, "+", "%20", -1)

	res, err := c.doRequest("GET", fmt.Sprintf("/api/orgs/name/%s", name), nil)
	if err != nil {
		return org, err
	}

	if !res.OK() {
		return org, res.Error()
	}

	err = res.BindJSON(&org)
	return org, err
}

// Orgs returns all the orgs in Grafana
func (c *Client) Orgs() ([]Org, error) {
	orgs := make([]Org, 0)

	res, err := c.doRequest("GET", "/api/orgs/", nil)
	if err != nil {
		return orgs, err
	}

	if !res.OK() {
		return orgs, res.Error()
	}

	err = res.BindJSON(&orgs)
	return orgs, err
}

// NewOrg creates an Org with the given name in Grafana
func (c *Client) NewOrg(name string) (Org, error) {
	org := Org{Name: name}
	newOrg := map[string]string{"name": name}
	res, err := c.doJSONRequest("POST", "/api/orgs", newOrg)
	if err != nil {
		return org, err
	}

	if !res.OK() {
		return org, res.Error()
	}

	body := struct {
		ID int64 `json:"orgId"`
	}{0}

	err = res.BindJSON(&body)
	if err == nil {
		org.ID = body.ID
	}

	return org, err
}

// DeleteOrg deletes the given org ID from Grafana
func (c *Client) DeleteOrg(id int64) error {
	res, err := c.doRequest("DELETE", fmt.Sprintf("/api/orgs/%d", id), nil)
	if err != nil {
		return err
	}

	return res.Error()
}

// UserRoleIsValid will return true if the given role is valid
func UserRoleIsValid(role string) bool {
	switch role {
	case OrgUserRoleAdmin:
		fallthrough
	case OrgUserRoleEditor:
		fallthrough
	case OrgUserRoleViewer:
		return true
	}

	return false
}

func AutoFixRole(role string) string {
	role = strings.Title(role)

	if role == "Veiwer" {
		role = OrgUserRoleViewer
	}

	return role
}
