package gapi

import (
	"fmt"
)

// AlertNotification represents a Grafana alert notification
type AlertNotification struct {
	Id        int64       `json:"id,omitempty"`
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	IsDefault bool        `json:"isDefault"`
	Settings  interface{} `json:"settings"`
}

// AlertNotification gets the alert with the given ID from Grafana
func (c *Client) AlertNotification(id int64) (*AlertNotification, error) {
	result := &AlertNotification{}
	path := fmt.Sprintf("/api/alert-notifications/%d", id)
	res, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if !res.OK() {
		return result, res.Error()
	}

	err = res.BindJSON(&result)
	return result, err
}

// NewAlertNotification creates the given alert notification object in Grafana
func (c *Client) NewAlertNotification(a *AlertNotification) (int64, error) {
	res, err := c.doJSONRequest("POST", "/api/alert-notifications", a)
	if err != nil {
		return 0, err
	}

	if !res.OK() {
		return 0, res.Error()
	}

	result := struct {
		ID int64 `json:"id"`
	}{}
	err = res.BindJSON(&result)
	return result.ID, err
}

// UpdateAlertNotification wll update the alert notification in Grafana that matches
// the ID from the given alert notification object
func (c *Client) UpdateAlertNotification(a *AlertNotification) error {
	path := fmt.Sprintf("/api/alert-notifications/%d", a.Id)
	res, err := c.doJSONRequest("PUT", path, a)
	if err != nil {
		return err
	}

	return res.Error()
}

// DeleteAlertNotification will delete the alert notification from Grafana
// matching the given ID
func (c *Client) DeleteAlertNotification(id int64) error {
	path := fmt.Sprintf("/api/alert-notifications/%d", id)
	res, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return res.Error()
}
