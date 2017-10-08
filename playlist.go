package gapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type Playlist struct {
	Id       int64          `json:"id,omitempty"`
	Name     string         `json:"name"`
	Interval string         `json:"interval"`
	URL      string         `json:"url"`
	Items    []PlaylistItem `json:"items"`
}

type PlaylistItem struct {
	Id    int64  `json:"id,omitempty"`
	Type  string `json:"type"`
	Order int    `json:"order"`
	Title string `json:"title"`
	Value string `json:"value"`
}

func (c *Client) NewPlaylist(s *Playlist) (int64, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return 0, err
	}
	req, err := c.newRequest("POST", "/api/playlists", bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, errors.New(resp.Status)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	result := struct {
		Id int64 `json:"id"`
	}{}
	err = json.Unmarshal(data, &result)
	return result.Id, err
}

func (c *Client) UpdatePlaylist(s *Playlist) error {
	path := fmt.Sprintf("/api/playlists/%d", s.Id)
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	req, err := c.newRequest("PUT", path, bytes.NewBuffer(data))
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

	return nil
}

func (c *Client) Playlist(id int64) (*Playlist, error) {
	path := fmt.Sprintf("/api/playlists/%d", id)
	req, err := c.newRequest("GET", path, nil)
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

	result := &Playlist{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (c *Client) DeletePlaylist(id int64) error {
	path := fmt.Sprintf("/api/playlists/%d", id)
	req, err := c.newRequest("DELETE", path, nil)
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

	return nil
}
