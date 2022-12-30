package entities

type Equipment struct {
	Key               string             `json:"key"`
	Name              string             `json:"name"`
	EquipmentCategory *EquipmentCategory `json:"equipment_category"`
	Cost              *Cost              `json:"cost"`
	Weight            int                `json:"weight"`
}

func (e *Equipment) GetType() string {
	return "equipment"
}

type EquipmentCategory struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Cost struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

type Weapon struct {
	Key               string             `json:"key"`
	Name              string             `json:"name"`
	EquipmentCategory *EquipmentCategory `json:"equipment_category"`
	Cost              *Cost              `json:"cost"`
	Weight            int                `json:"weight"`
	WeaponCategory    string             `json:"weapon_category"`
	Damage            *Damage            `json:"damage"`
	WeaponRange       string             `json:"range"`
	CategoryRange     string             `json:"category_range"`
	Range             *Range             `json:"weapon_range"`
	Properties        []*Properties      `json:"properties"`
	TwoHandedDamage   *Damage            `json:"two_handed_damage"`
}

func (w *Weapon) GetType() string {
	return "weapon"
}

type Damage struct {
	DamageDice string      `json:"damage_dice"`
	DamageType *DamageType `json:"damage_type"`
}

type DamageType struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Properties struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Range struct {
	Normal int `json:"normal"`
}
