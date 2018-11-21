package gapi

import (
	"testing"

	"github.com/gobs/pretty"
)

const (
	creatTeamJSON = `{"message":"Team created","teamId":2}`
	getTeamJSON   = `{
                      "id": 1,
                      "orgId": 1,
                      "name": "MyTestTeam",
                      "email": "",
                      "created": "2017-12-15T10:40:45+01:00",
                      "updated": "2017-12-15T10:40:45+01:00"
                    }`
	updateTeamJSON    = `{"message":"Team updated"}`
	deleteTeamJSON    = `{"message":"Team deleted"}`
	addTeamMemberJSON = `{"message":"Member added to Team"}`
)

func TestNewTeam(t *testing.T) {
	server, client := gapiTestTools(200, creatTeamJSON)
	defer server.Close()

	team, err := client.NewTeam("test")
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(team))

	if team.Id != 2 {
		t.Error("Create Team Error.")
	}
}

func TestGetTeam(t *testing.T) {
	server, client := gapiTestTools(200, getTeamJSON)
	defer server.Close()

	resp, err := client.Team(1)
	if err != nil {
		t.Error(err)
	}

	t.Log(pretty.PrettyFormat(resp))

	team := Team{
		Id:    1,
		OrgId: 1,
		Email: "",
		Name:  "MyTestTeam",
	}
	if resp != team {
		t.Error("Not correctly parsing returned team.")
	}
}

func TestUpdateTeam(t *testing.T) {
	server, client := gapiTestTools(200, updateTeamJSON)
	defer server.Close()

	err := client.UpdateTeam("1", "test team 2")
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteTeam(t *testing.T) {
	server, client := gapiTestTools(200, deleteTeamJSON)
	defer server.Close()

	err := client.DeleteTeam("1")
	if err != nil {
		t.Error(err)
	}
}

func TestAddTeamMember(t *testing.T) {
	server, client := gapiTestTools(200, addTeamMemberJSON)
	defer server.Close()

	err := client.AddTeamMember("1", 1)
	if err != nil {
		t.Error(err)
	}
}
