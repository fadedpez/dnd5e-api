package entities

type Equipment struct {
	Key               string         `json:"key"`
	Name              string         `json:"name"`
	EquipmentCategory *ReferenceItem `json:"equipment_category"`
	Cost              *Cost          `json:"cost"`
	Weight            float32        `json:"weight"`
}

func (e *Equipment) GetType() string {
	return "equipment"
}

type Cost struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

type Weapon struct {
	Key               string           `json:"key"`
	Name              string           `json:"name"`
	EquipmentCategory *ReferenceItem   `json:"equipment_category"`
	Cost              *Cost            `json:"cost"`
	Weight            float32          `json:"weight"`
	WeaponCategory    string           `json:"weapon_category"`
	Damage            *Damage          `json:"damage"`
	WeaponRange       string           `json:"range"`
	CategoryRange     string           `json:"category_range"`
	Range             *Range           `json:"weapon_range"`
	Properties        []*ReferenceItem `json:"properties"`
	TwoHandedDamage   *Damage          `json:"two_handed_damage"`
}

func (w *Weapon) GetType() string {
	return "weapon"
}

type Damage struct {
	DamageDice string         `json:"damage_dice"`
	DamageType *ReferenceItem `json:"damage_type"`
}

type Range struct {
	Normal int `json:"normal"`
}

type Armor struct {
	Key                 string         `json:"key"`
	Name                string         `json:"name"`
	EquipmentCategory   *ReferenceItem `json:"equipment_category"`
	Cost                *Cost          `json:"cost"`
	Weight              float32        `json:"weight"`
	ArmorCategory       string         `json:"armor_category"`
	ArmorClass          *ArmorClass    `json:"armor_class"`
	StrMinimum          int            `json:"str_minimum"`
	StealthDisadvantage bool           `json:"stealth_disadvantage"`
}

func (a *Armor) GetType() string {
	return "armor"
}

type ArmorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
}
