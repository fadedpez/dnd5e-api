package entities

type Class struct {
	Key               string               `json:"key"`
	Name              string               `json:"name"`
	HitDie            int                  `json:"hit_die"`
	Proficiencies     []*Proficiency       `json:"proficiencies"`
	SavingThrows      []*SavingThrow       `json:"saving_throws"`
	StartingEquipment []*StartingEquipment `json:"starting_equipment"`
}

type SavingThrow struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type StartingEquipment struct {
	Equipment *EquipmentList `json:"equipment"`
	Quantity  int            `json:"quantity"`
}

type EquipmentList struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
