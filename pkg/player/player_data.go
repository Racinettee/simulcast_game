package player

import "github.com/Racinettee/simul/pkg/component"

// This data is for the parts of the player that should persist across scenes
// eg, current health, current wealth, etc and inventory which can all be serialized into json
// this is no place for information about current animation frame
type PlayerDat struct {
	CurrentHealth int                       `json:"current_health"`
	CurrentWealth int                       `json:"current_wealth"`
	Inventory     map[string]component.Item `json:"inventory"`
	EquipedWeapon component.Weapon
}
