package gapi

import (
	"bytes"
	"encoding/json"

	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)


type FolderCreateResponse struct {
	Id     int `json:"id"`
	uid    string `json:"uid"`
	title  string `json:"title"`
}

type Folder struct {
	//Model map[string]interface{} `json:"folder"`
	Id int `json:"id"`
}

func (c *Client) CreateFolder(model map[string]interface{}) (*FolderCreateResponse, error) {
	wrapper := map[string]interface{}{
		"title": model["title"],
		"uid": 	 model["uid"],
	}
	data, err := json.Marshal(wrapper)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshall folder JSON")
	}
	req, err := c.newRequest("POST", "/api/folders", nil, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to perform HTTP request")
	}

	if resp.StatusCode != 200 {
		var gmsg GrafanaErrorMessage
		dec := json.NewDecoder(resp.Body)
		dec.Decode(&gmsg)
		return nil, fmt.Errorf("Request to Grafana returned %+v status code with the following message: %+v", resp.StatusCode, gmsg.Message)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &FolderCreateResponse{}
	err = json.Unmarshal(data, &result)
	return result, err
}

func (c *Client) Folder(slug string) (*Folder, error) {
	path := fmt.Sprintf("/api/folders/%s", slug)
	req, err := c.newRequest("GET", path, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	// If the error code is 404 that means that the folder does not exist. Don't treat this case as an error
	if resp.StatusCode == 404 {
		return nil, nil
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &Folder{}
	err = json.Unmarshal(data, &result)
	return result, err
}

