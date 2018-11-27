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

type TeamMember struct {
	UserId    int64  `json:"userId,omitempty"`
	TeamId    int64  `json:"teamId,omitempty"`
	OrgId     int64  `json:"orgId,omitempty"`
	Login     string `json:"login,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarUrl string `json:"avatarUrl,omitempty"`
}

type TeamSearchResponse struct {
	TotalCount int64   `json:"totalCount,omitempty"`
	Teams      []*Team `json:"teams,omitempty"`
	Page       int64   `json:"page,omitempty"`
	PerPage    int64   `json:"perPage,omitempty"`
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

func (c *Client) Teams() ([]*Team, error) {
	var list []*Team
	req, err := c.newRequest("GET", "/api/teams/search", nil, nil)
	if err != nil {
		return list, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return list, err
	}
	if resp.StatusCode != 200 {
		return list, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return list, err
	}
	var r TeamSearchResponse
	err = json.Unmarshal(data, &r)
	return r.Teams, err
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

func (c *Client) AddTeamMember(id string, userID int64) error {
	dataMap := map[string]interface{}{
		"userId": userID,
	}
	data, err := json.Marshal(dataMap)
	if err != nil {
		return err
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/api/teams/%s/members", id), nil, bytes.NewBuffer(data))
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

func (c *Client) RemoveTeamMember(id string, userID int64) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/api/teams/%s/members/%d", id, userID), nil, nil)
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

func (c *Client) TeamMembers(id string) ([]*TeamMember, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/api/teams/%s/members", id), nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var list []*TeamMember
	err = json.Unmarshal(data, &list)
	return list, err
}
