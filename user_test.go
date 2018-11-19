package gapi

import (
	"testing"

	"github.com/gobs/pretty"
)

const (
	getUsersJSON       = `[{"id":1,"name":"","login":"admin","email":"admin@localhost","avatarUrl":"/avatar/46d229b033af06a191ff2267bca9ae56","isAdmin":true,"lastSeenAt":"2018-06-28T14:42:24Z","lastSeenAtAge":"\u003c 1m"}]`
	getUserJSON        = `{"id": 1, "email": "user@mygraf.com", "name": "admin", "login": "admin", "theme": "light", "orgId": 1, "isGrafanaAdmin": true}`
	getUserByEmailJSON = `{"id":1,"email":"admin@localhost","name":"","login":"admin","theme":"","orgId":1,"isGrafanaAdmin":true}`
)

func TestUsers(t *testing.T) {
	server, client := gapiTestTools(200, getUsersJSON)
	defer server.Close()

	resp, err := client.Users()
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(resp))

	user := User{
		Id:      1,
		Email:   "admin@localhost",
		Name:    "",
		Login:   "admin",
		IsAdmin: true,
	}

	if len(resp) != 1 || resp[0] != user {
		t.Error("Not correctly parsing returned users.")
	}
}

func TestGetUser(t *testing.T) {
	server, client := gapiTestTools(200, getUserJSON)
	defer server.Close()

	resp, err := client.User(1)
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(resp))

	user := User{
		Id:      1,
		Email:   "user@mygraf.com",
		Name:    "admin",
		Login:   "admin",
		IsAdmin: true,
	}
	if resp.Id != user.Id {
		t.Error("Not correctly parsing returned user.")
	}
}

func TestUserByEmail(t *testing.T) {
	server, client := gapiTestTools(200, getUserByEmailJSON)
	defer server.Close()

	resp, err := client.UserByEmail("admin@localhost")
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(resp))

	user := User{
		Id:      1,
		Email:   "admin@localhost",
		Name:    "",
		Login:   "admin",
		IsAdmin: true,
	}
	if resp != user {
		t.Error("Not correctly parsing returned user.")
	}
}
