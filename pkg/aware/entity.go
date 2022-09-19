package aware

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

// GetParentHierachyName returns a string of the full entity hierachy path.
// i.e. AGL Conveyor Motor Sensor instead of just Sensor.
func (e *Entity) GetParentHierachyName() string {
	if e.ParentEntity == nil {
		return e.Name
	}
	return e.ParentEntity.GetParentHierachyName() + " " + e.Name
}
