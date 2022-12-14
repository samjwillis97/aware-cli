package aware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-faker/faker/v4"
)

// DeviceType is the aware model for a Device Type.
type DeviceType struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Kind         string      `json:"kind"`
	Options      interface{} `json:"options"`
	Description  string      `json:"description"`
	IsShared     bool        `json:"isShared"`
	IsActive     bool        `json:"isActive"`
	IsHidden     bool        `json:"isHidden"`
	Organisation string      `json:"organisation"`
	Scope        string      `json:"scope"`
	// AllowedAttributes
	Parameters []DeviceTypeParameter `json:"parameters"`
	// DisplayGroups
	// Commands
}

// DeviceTypeParameter is the aware model for a parameter.
type DeviceTypeParameter struct {
	ID             string                       `json:"id"`
	Name           string                       `json:"name"`
	DisplayName    string                       `json:"displayName"`
	ValueType      DeviceTypeParameterValueType `json:"valueType"`
	IsActive       bool                         `json:"isActive"`
	Display        DeviceTypeParameterDisplay   `json:"display"`
	IsPrimary      bool                         `json:"isPrimary"`
	IsAggregatable bool                         `json:"IsAggregatable"`
	Range          DeviceTypeParameterRange     `json:"range"`
	// Alarms
}

// DeviceTypeParameterRange is the aware model for parameter ranges.
type DeviceTypeParameterRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// DeviceTypeParameterDisplay is the aware model for parameter displays.
type DeviceTypeParameterDisplay struct {
	Unit      string                   `json:"unit"`
	Scale     float64                  `json:"scale"`
	Range     DeviceTypeParameterRange `json:"range"`
	Component string                   `json:"component"`
	// Values
}

// DeviceTypeParameterValueType is the aware enum for parameter value types.
type DeviceTypeParameterValueType string

const (
	// Float matches the aware parameter value type.
	Float DeviceTypeParameterValueType = "float"
	// Bool matches the aware parameter value type.
	Bool DeviceTypeParameterValueType = "bool"
	// String matches the aware parameter value type.
	String DeviceTypeParameterValueType = "string"
	// Object matches the aware parameter value type.
	Object DeviceTypeParameterValueType = "object"
	// Waveform matches the aware parameter value type.
	Waveform DeviceTypeParameterValueType = "waveform"
	// Spectrum matches the aware parameter value type.
	Spectrum DeviceTypeParameterValueType = "spectrum"
)

// GetAllDeviceTypes attempts to retrieve all device types.
func (c *Client) GetAllDeviceTypes(org string) ([]*DeviceType, error) {
	url := fmt.Sprintf("%s/v1/devicetypes", c.server)

	if org != "" {
		url += "?organisationId=" + org
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

	var out []*DeviceType
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

// GetDeviceTypeByID attempts to get the device with the given ID.
func (c *Client) GetDeviceTypeByID(id string) (*DeviceType, error) {
	url := fmt.Sprintf("%s/v1/devicetypes/%s", c.server, id)

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

	var out *DeviceType
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

// GetRandomValue generates a random value for the parameter.
// nolint:gocyclo // Complexity is required to generate more realistic random values
func (p *DeviceTypeParameter) GetRandomValue() interface{} {
	// Copied staight from Jez's LinqPad
	// FIXME: Doesn't seem to be respecting these
	if (p.Display != DeviceTypeParameterDisplay{} && p.Display.Unit != "") {
		switch p.Display.Unit {
		case "ohm", "resistance":
			return generateRandomFloat(0, 60, 2)
		case "volt-ampere", "volt-ampere-reactive":
			return generateRandomFloat(0, 3000, 2)
		case "watt-hour":
			return generateRandomFloat(1000, 50000000, 2)
		case "watt":
			if strings.Contains(p.DisplayName, "DC") {
				return generateRandomFloat(1, 90, 2)
			}
			return generateRandomFloat(1, 1500, 2)
		case "amp", "amps", "ampere":
			return generateRandomFloat(1, 100, 2)
		case "volt", "volts", "voltage":
			return generateRandomFloat(200, 300, 2)
		case "percent", "percentage":
			return generateRandomFloat(1, 100, 2)
		case "hertz", "frequency":
			return generateRandomFloat(1, 60, 2)
		case "degrees-celsius":
			return generateRandomFloat(20, 100, 2)
		}
	}

	name := strings.ToLower(p.DisplayName)
	switch {
	case strings.Contains(name, "running"):
		return true
	case strings.Contains(name, "voltage"):
		return generateRandomFloat(200, 300, 2)
	case strings.Contains(name, "setting"):
		return generateRandomInt(0, 5)
	case strings.Contains(name, "factor"):
		return generateRandomFloat(0.1, 0.99, 2)
	case strings.Contains(name, "test-report"):
		return nil
	case strings.Contains(name, "status"):
		return generateRandomBool()
	case strings.Contains(name, "date"):
		// TODO: Now
		// faker.Date()
		return nil
	case strings.Contains(name, "speed"):
		return generateRandomInt(0, 10)
	case strings.Contains(name, "frequency"):
		return generateRandomFloat(1, 60, 2)
	case strings.Contains(name, "power"):
		return generateRandomFloat(0, 30, 2)
	}

	switch p.ValueType {
	case Float:
		return generateRandomFloat(0, 100, 2)
	case Bool:
		return generateRandomBool()
	case String:
		return faker.Word()
	}

	return nil
}
