package aware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Device struct {
	ID           string     `json:"id"`
	DeviceType   DeviceType `json:"deviceType"`
	IsActive     bool       `json:"isActive"`
	IsEnabled    bool       `json:"isEnabled"`
	IsHidden     bool       `json:"isHidden"`
	ParentEntity Entity     `json:"parentEntity"`
	Organisation string     `json:"organisation"`
	CloudID      string     `json:"cloudId"`
	// Attributes
	// Identity
	// IdentityHistory
	// Credentials
	// LatestValues
	DisplayName string      `json:"displayName"`
	State       interface{} `json:"state"`
}

type GetAllDevicesOptions struct {
	IncludeInactive     bool
	EntityId            string
	OrganisationId      string
	DeviceTypeKind      string
	IncludeLatestValues bool
}

func (c *Client) GetAllDevices(opts GetAllDevicesOptions) ([]*Device, error) {
	queryString := ""

	// TODO: Test
	if opts.DeviceTypeKind != "" {
		queryString += fmt.Sprintf("deviceTypeKind=%s", opts.DeviceTypeKind)
	}
	if opts.IncludeInactive {
		if queryString != "" {
			queryString += "&"
		}
		queryString += fmt.Sprintf("includeInactive=%v", opts.IncludeInactive)
	}
	if opts.EntityId != "" {
		if queryString != "" {
			queryString += "&"
		}
		queryString += fmt.Sprintf("entityId=%s", opts.EntityId)
	}
	if opts.OrganisationId != "" {
		if queryString != "" {
			queryString += "&"
		}
		queryString += fmt.Sprintf("organisationId=%s", opts.OrganisationId)
	}
	if opts.IncludeLatestValues {
		if queryString != "" {
			queryString += "&"
		}
		queryString += fmt.Sprintf("IncludeLatestValues=%v", opts.IncludeLatestValues)
	}

	url := c.server + "/v1/devices"
	if queryString != "" {
		url += "?" + queryString
	}

	res, err := c.request(context.Background(), http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, ErrEmptyResult
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		// TODO: Pretty Print?
		return nil, err
	}

	var out []*Device
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) GetDeviceByID(id string) (*Device, error) {
	url := c.server + "/v1/devices/" + id

	res, err := c.request(context.Background(), http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, ErrEmptyResult
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var out *Device
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}
