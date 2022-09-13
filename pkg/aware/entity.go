package aware

type Entity struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Description string `json:"description"`
    // Status
    // Attributes
    ParentEntity *Entity `json:"parentEntity"`
    Organisation string `json:"organisation"`
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

// TODO: Need to actually test this with more than one parent
func (e *Entity) GetParentHierachyName() string {
    if e.ParentEntity == nil {
        return e.Name
    }
    return e.ParentEntity.GetParentHierachyName() + " " + e.Name
}
