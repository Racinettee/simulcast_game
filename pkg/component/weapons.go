package component

//go:generate stringer -type=WeaponType

type WeaponType byte

const (
	Sword WeaponType = iota
	Spear
	Flail
	Mace
	Bow
	Staff
	Fist
)

type Weapon struct {
	Name       string
	Type       WeaponType
	BaseDamage byte
}

var IronSword = Weapon{
	Name:       "IronSword",
	Type:       Sword,
	BaseDamage: 1,
}

var IronSpear = Weapon{
	Name:       "IronSpear",
	Type:       Spear,
	BaseDamage: 2,
}
