package aware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Device is the aware model of a device.
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

// GetAllDevicesOptions are the available options for the GetAllDevices query.
type GetAllDevicesOptions struct {
	IncludeInactive     bool
	EntityID            string
	OrganisationID      string
	DeviceTypeKind      string
	IncludeLatestValues bool
}

// GetAllDevices gets all the available devices for a user with the given options.
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
	if opts.EntityID != "" {
		if queryString != "" {
			queryString += "&"
		}
		queryString += fmt.Sprintf("entityId=%s", opts.EntityID)
	}
	if opts.OrganisationID != "" {
		if queryString != "" {
			queryString += "&"
		}
		queryString += fmt.Sprintf("organisationId=%s", opts.OrganisationID)
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

// GetDeviceByID attempts to retrieve a device with the given id.
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
