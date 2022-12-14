package aware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Device is the aware model of a device when listing.
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

// CreatedDevice is the aware model return when creating a device.
type CreatedDevice struct {
	ID           string `json:"id"`
	DeviceType   string `json:"deviceType"`
	IsActive     bool   `json:"isActive"`
	IsEnabled    bool   `json:"isEnabled"`
	IsHidden     bool   `json:"isHidden"`
	ParentEntity string `json:"parentEntity"`
	Organisation string `json:"organisation"`
	CloudID      string `json:"cloudId"`
	// Attributes
	// Identity
	// IdentityHistory
	// Credentials
	// LatestValues
	DisplayName string `json:"displayName"`
}

// GetAllDevicesOptions are the available options for the GetAllDevices query.
type GetAllDevicesOptions struct {
	IncludeInactive     bool
	EntityID            string
	OrganisationID      string
	DeviceTypeKind      string
	IncludeLatestValues bool
}

// CreateDeviceRequest is the data used to create a new device.
type CreateDeviceRequest struct {
	DisplayName  string `json:"displayName"`
	DeviceType   string `json:"deviceType"`
	ParentEntity string `json:"parentEntity"`
	Organisation string `json:"organisation"`
	// IsActive bool
	// IsEnabled bool
	// Identity
	// IdentityHistory
	// Credentials
}

// UpdateDeviceRequest is the data used when updating an existing device.
type UpdateDeviceRequest struct {
	DeviceType   string `json:"deviceType"`
	ParentEntity string `json:"parentEntity"`
	Organisation string `json:"organisation"`
	DisplayName  string `json:"displayName"`
}

// CreateDevice will create a new device with the given request details.
func (c *Client) CreateDevice(req *CreateDeviceRequest) (*CreatedDevice, error) {
	header := Header{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	body, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	res, err := c.request(context.Background(), http.MethodPost, c.server+"/v1/devices", body, header)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, ErrEmptyResult
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusCreated {
		return nil, formatUnexpectedResponse(res)
	}

	var out *CreatedDevice
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

// DeleteDevice will delete a device with the given ID.
func (c *Client) DeleteDevice(id string) error {
	res, err := c.request(context.Background(), http.MethodDelete, c.server+"/v1/devices/delete/"+id, nil, nil)
	if err != nil {
		return err
	}

	if res == nil {
		return ErrEmptyResult
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNoContent {
		return formatUnexpectedResponse(res)
	}

	return nil
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
		return nil, formatUnexpectedResponse(res)
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

// UpdateDeviceByID updates details of a device on aware.
func (c *Client) UpdateDeviceByID(id string, req *UpdateDeviceRequest) error {
	header := Header{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	body, err := json.Marshal(&req)
	if err != nil {
		return err
	}

	res, err := c.request(context.Background(), http.MethodPut, c.server+"/v1/devices/update/"+id, body, header)
	if err != nil {
		return err
	}

	if res == nil {
		return ErrEmptyResult
	}

	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNoContent {
		return formatUnexpectedResponse(res)
	}

	return nil
}
