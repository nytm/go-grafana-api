package gapi

import (
	"fmt"
)

// DataSource represents a Grafana data source
type DataSource struct {
	ID     int64  `json:"id,omitempty"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	URL    string `json:"url"`
	Access string `json:"access"`

	Database string `json:"database,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`

	OrgID     int64 `json:"orgId,omitempty"`
	IsDefault bool  `json:"isDefault"`

	BasicAuth         bool   `json:"basicAuth"`
	BasicAuthUser     string `json:"basicAuthUser,omitempty"`
	BasicAuthPassword string `json:"basicAuthPassword,omitempty"`

	JSONData       JSONData       `json:"jsonData,omitempty"`
	SecureJSONData SecureJSONData `json:"secureJsonData,omitempty"`
}

// JSONData is a representation of the datasource `jsonData` property
type JSONData struct {
	AssumeRoleArn string `json:"assumeRoleArn,omitempty"`
	AuthType      string `json:"authType,omitempty"`
	DefaultRegion string `json:"defaultRegion,omitempty"`
}

// SecureJSONData is a representation of the datasource `secureJsonData` property
type SecureJSONData struct {
	AccessKey string `json:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty"`
}

// NewDataSource will create the given data source in Grafana
func (c *Client) NewDataSource(s *DataSource) (int64, error) {
	res, err := c.doJSONRequest("POST", "/api/datasources", s)
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

// UpdateDataSource will update the data source in Grafana from the given
// datasource object that matches the given datasource objects ID
func (c *Client) UpdateDataSource(s *DataSource) error {
	path := fmt.Sprintf("/api/datasources/%d", s.ID)
	res, err := c.doJSONRequest("PUT", path, s)
	if err != nil {
		return err
	}

	return res.Error()
}

// DataSource will return the datasource with the given ID
func (c *Client) DataSource(id int64) (*DataSource, error) {
	path := fmt.Sprintf("/api/datasources/%d", id)
	res, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	if !res.OK() {
		return nil, res.Error()
	}

	result := &DataSource{}
	err = res.BindJSON(&result)
	return result, err
}

// DeleteDataSource will delete the datasource with the given ID from Grafana
func (c *Client) DeleteDataSource(id int64) error {
	path := fmt.Sprintf("/api/datasources/%d", id)
	res, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return res.Error()
}

// DataSourcesByOrgID will return the datasources for the given org ID
func (c *Client) DataSourcesByOrgID(id int64) ([]*DataSource, error) {
	out := []*DataSource{}
	dss, err := c.DataSources()
	if err != nil {
		return out, err
	}

	for _, ds := range dss {
		if ds.OrgID == id {
			out = append(out, ds)
		}
	}

	return out, nil
}

// DataSources will return all the datasources from Grafana
func (c *Client) DataSources() ([]*DataSource, error) {
	path := "/api/datasources"
	res, err := c.doRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	result := []*DataSource{}
	err = res.BindJSON(&result)
	return result, err
}
