package aware

import (
	"context"
	"encoding/json"
	"net/http"
)

type AuthProvider struct {
	AuthType         string      `json:"type"`
	Label            string      `json:"label"`
	IsDefault        bool        `json:"isDefault"`
	IsExternal       bool        `json:"isExternal"`
	CanResetPassword bool        `json:"canResetPassword"`
	IsHidden         bool        `json:"isHidden"`
	UseForm          bool        `json:"useForm"`
	Url              string      `json:"url"`
	Data             interface{} `json:"data"`
}

type AuthResponse struct {
	ID          string `json:"id"`
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
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

func (c *Client) Login(login, password, providerType string) (*AuthResponse, error) {
    // TODO: Maybe need to change the providerType to feed in the full
    // AuthProvider so we have the URL etc, for Auth0
	data := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Provider string `json:"provider"`
	}{login, password, providerType}

	header := Header{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	body, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	res, err := c.request(context.Background(), http.MethodPost, c.server+"/v1/auth/login", body, header)
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

	var out *AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}
