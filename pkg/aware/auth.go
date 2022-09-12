package aware

import (
	"context"
	"encoding/json"
	"net/http"
)

type AuthProvider struct {
    AuthType string `json:"type"`
    Label string `json:"label"`
    IsDefault bool `json:"isDefault"`
    IsExternal bool `json:"isExternal"`
    CanResetPassword bool `json:"canResetPassword"`
    IsHidden bool `json:"isHidden"`
    UseForm bool `json:"useForm"`
    Url string `json:"url"`
    Data interface{} `json:"data"`
}

func (c *Client) GetAllAuthProviders() ([]*AuthProvider, error) {
    res, err := c.request(context.Background(), http.MethodGet, c.server+"/v1/auth/providers", nil, nil) 
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

    var out []*AuthProvider
    if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
        return nil, err
    }

    return out, nil
}
