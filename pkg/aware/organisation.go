package aware

import (
    "net/http"
    "context"
    "encoding/json"
)

type Organisation struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
    Options interface{} `json:"options"`
    Abbreviation string `json:"abbreviation"`
    AllowedDeviceTypes []string `json:"allowedDeviceTypes"`
    AllowedEntityTypes []string `json:"allowedEntityTypes"`
    AllowedActivityTypes []string `json:"allowedActivityTypes"`
    Files []OrganisationFile `json:"files"`
}

type OrganisationFile struct {
    ID  string `json:"id"`
    Name string `json:"name"`
    Type string `json:"string"`
    URI interface {} `json:"uri"`
}

func (c *Client) GetAllOrganisations() ([]*Organisation, error) {
	res, err := c.request(context.Background(), http.MethodGet, c.server+"/v1/organisations", nil, nil)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, ErrEmptyResult
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		// TODO: Pretty Print
		return nil, err
	}

	var out []*Organisation
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}
