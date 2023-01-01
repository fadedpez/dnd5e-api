package entities

type Class struct {
	Key               string               `json:"key"`
	Name              string               `json:"name"`
	HitDie            int                  `json:"hit_die"`
	Proficiencies     []*ReferenceItem       `json:"proficiencies"`
	SavingThrows      []*ReferenceItem       `json:"saving_throws"`
	StartingEquipment []*StartingEquipment `json:"starting_equipment"`
}

type StartingEquipment struct {
	Equipment *ReferenceItem `json:"equipment"`
	Quantity  int            `json:"quantity"`
}

