package aware

type DeviceType struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Kind string `json:"kind"`
    Options interface{} `json:"options"`
    Description string `json:"description"`
    IsShared bool `json:"isShared"`
    IsActive bool `json:"isActive"`
    IsHidden bool `json:"isHidden"`
    Organisation string `json:"organisation"`
    Scope string `json:"scope"`
    // AllowedAttributes
    // Parameters
    // DisplayGroups
    // Commands
}
