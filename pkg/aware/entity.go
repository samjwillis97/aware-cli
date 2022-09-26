package aware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Entity is the aware model an entity.
type Entity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// Status
	// Attributes
	ParentEntity *Entity `json:"parentEntity"`
	Organisation string  `json:"organisation"`
	// EntityType
	IsActive bool `json:"isActive"`
	// Gateways
	// Files
	// Path
	// Ancestors
	// Note
	// Order
	// Criticality
	// ComputedCriticality
	// InheritedCriticality
	// Identity
	// IdentityHistory
}

// GetAllEntitiesOptions are the available options for the GetAllEntities query.
type GetAllEntitiesOptions struct {
	IncludeInactive bool
	ExcludeDetail   bool
	ParentEntityID  string
	Group           string
	Kind            string
}

// GetAllEntities attempts to retrieve all entities for an organistion.
// org is required.
func (c *Client) GetAllEntities(org string, opts GetAllEntitiesOptions) ([]*Entity, error) {
	// TODO: Opts
	url := fmt.Sprintf("%s/v1/entities?organisationId=%s", c.server, org)

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

	var out []*Entity
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

// GetParentHierachyName returns a string of the full entity hierachy path.
// i.e. AGL Conveyor Motor Sensor instead of just Sensor.
func (e *Entity) GetParentHierachyName() string {
	if e.ParentEntity == nil {
		return e.Name
	}
	return e.ParentEntity.GetParentHierachyName() + " - " + e.Name
}
