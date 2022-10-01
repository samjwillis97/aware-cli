package aware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestGetDeviceByID(t *testing.T) {
	var unexpectedStatusCode bool

	is := is.New(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		is.Equal("/v1/devices/TEST-1", r.URL.Path)

		if unexpectedStatusCode {
			w.WriteHeader(400)
		} else {
			resp, err := os.ReadFile("./test_data/device.json")
			is.NoErr(err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			_, _ = w.Write(resp)
		}
	}))
	defer server.Close()

	client := NewClient(Config{Server: server.URL})

	actual, err := client.GetDeviceByID("TEST-1")
	is.NoErr(err)

	expected := &Device{
		ID:        "5d1d574439d157849090ea6a",
		IsActive:  true,
		IsEnabled: true,
		IsHidden:  false,
		CloudID:   "AMP_HST_H266_Outlet_2_IPB_5d1d574439d157849090ea6a",
		DeviceType: DeviceType{
			IsShared:    false,
			IsActive:    true,
			IsHidden:    false,
			ID:          "5cf717e2bec882982729dd8a",
			Name:        "IPB",
			Kind:        "integrated-protection-relay",
			Description: "Integrated Protection Relay",
			Scope:       "system",
			// "allowedAttributes": [],
			Parameters: []DeviceTypeParameter{
				{
					ID:             "5d1af1d930cbe93bdcba8239",
					Name:           "pilot-forward-resistance",
					DisplayName:    "Pilot Forward Resistance",
					ValueType:      "float",
					IsActive:       true,
					IsPrimary:      false,
					IsAggregatable: false,
					Display: DeviceTypeParameterDisplay{
						Unit:  "ohm",
						Scale: 1,
					},
					Range: DeviceTypeParameterRange{
						Min: 0,
						Max: 150,
					},
				},
			},
			// "displayGroups": [
			//     {
			//         "id": "5d3142fb4ed8460012ea2628",
			//         "name": "Settings",
			//         "type": "list",
			//         "displayName": "Settings",
			//         "order": 10,
			//         "selectedParameters": [
			//         "5d1af1d930cbe93bdcba823b",
			//         "5d1af1d930cbe93bdcba823c",
			//         "5d1af1d930cbe93bdcba823d",
			//         "5d1af1d930cbe93bdcba823e",
			//         "5d1af1d930cbe93bdcba823f",
			//         "5d1af1d930cbe93bdcba8240",
			//         "5d1af1d930cbe93bdcba8251",
			//         "5d1af1d930cbe93bdcba8252"
			//         ]
			//     }
			// ]
		},
		ParentEntity: Entity{
			ID:          "5cf48c71b2f30979bc612292",
			Name:        "Outlet 2",
			Description: "DL032",
		},
		Organisation: "5bff4a241c7bed480ff3e261",
		// "identityHistory": []
		// ],
		// "attributes": []
	}
	is.Equal(expected, actual)

	unexpectedStatusCode = true
	_, err = client.GetDeviceByID("TEST-1")
	is.Equal(err, &ErrUnexpectedResponse{
		StatusCode: 400,
		Status:     "400 Bad Request",
	})
}
