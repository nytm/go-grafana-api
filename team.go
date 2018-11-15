package gapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Team struct {
	Id          int64  `json:"id,omitempty"`
	OrgId       int64  `json:"orgId,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	AvatarUrl   string `json:"avatarUrl,omitempty"`
	MemberCount int64  `json:"memberCount,omitempty"`
}

type CreateTeamResponse struct {
	Id int64 `json:"teamId"`
}

func (c *Client) Team(id int64) (Team, error) {
	team := Team{}
	req, err := c.newRequest("GET", fmt.Sprintf("/api/teams/%d", id), nil, nil)
	if err != nil {
		return team, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return team, err
	}
	if resp.StatusCode != 200 {
		return team, errors.New(resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return team, err
	}
	err = json.Unmarshal(data, &team)
	return team, err
}

func (c *Client) NewTeam(name string) (Team, error) {
	team := Team{
		Name: name,
	}
	data, err := json.Marshal(team)
	req, err := c.newRequest("POST", "/api/teams", nil, bytes.NewBuffer(data))
	if err != nil {
		return team, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return team, err
	}
	if resp.StatusCode != 200 {
		data, _ = ioutil.ReadAll(resp.Body)
		return team, fmt.Errorf("status: %s body: %s", resp.Status, data)
	}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return team, err
	}
	var r CreateTeamResponse
	err = json.Unmarshal(data, &r)
	if err != nil {
		return team, err
	}
	team.Id = r.Id
	return team, err
}

func (c *Client) UpdateTeam(id string, name string) error {
	dataMap := map[string]string{
		"name": name,
	}
	data, err := json.Marshal(dataMap)
	req, err := c.newRequest("PUT", fmt.Sprintf("/api/teams/%s", id), nil, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return err
}

func (c *Client) DeleteTeam(id string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/teams/%s", id), nil, nil)
	if err != nil {
		return err
	}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	return err
}
