package entities

type Class struct {
	Key                      string               `json:"key"`
	Name                     string               `json:"name"`
	HitDie                   int                  `json:"hit_die"`
	Proficiencies            []*ReferenceItem     `json:"proficiencies"`
	SavingThrows             []*ReferenceItem     `json:"saving_throws"`
	StartingEquipment        []*StartingEquipment `json:"starting_equipment"`
	ProficiencyChoices       []*ChoiceOption      `json:"proficiency_choices"`
	StartingEquipmentOptions []*ChoiceOption      `json:"starting_equipment_options"`
	PrimaryAbilities         []*ReferenceItem     `json:"primary_abilities"`
	Description              string               `json:"description"`
	ArmorProficiencies       []*ReferenceItem     `json:"armor_proficiencies"`
	WeaponProficiencies      []*ReferenceItem     `json:"weapon_proficiencies"`
	ToolProficiencies        []*ReferenceItem     `json:"tool_proficiencies"`
	Spellcasting             *ClassSpellcasting   `json:"spellcasting"`
}

type StartingEquipment struct {
	Equipment *ReferenceItem `json:"equipment"`
	Quantity  int            `json:"quantity"`
}

type ClassSpellcasting struct {
	Level               int              `json:"level"`
	SpellcastingAbility *ReferenceItem   `json:"spellcasting_ability"`
	Info                []*SpellcastingInfo `json:"info"`
}

type SpellcastingInfo struct {
	Name string   `json:"name"`
	Desc []string `json:"desc"`
}
