package entities

import "github.com/fadedpez/dnd5e-api/entities/choice"

type Class struct {
	Key                      string                 `json:"key"`
	Name                     string                 `json:"name"`
	HitDie                   int                    `json:"hit_die"`
	Proficiencies            []*ReferenceItem       `json:"proficiencies"`
	SavingThrows             []*ReferenceItem       `json:"saving_throws"`
	StartingEquipment        []*StartingEquipment   `json:"starting_equipment"`
	ProficiencyChoices       []*choice.ChoiceOption `json:"proficiency_choices"`
	StartingEquipmentOptions []*choice.ChoiceOption `json:"starting_equipment_options"`
}

type StartingEquipment struct {
	Equipment *ReferenceItem `json:"equipment"`
	Quantity  int            `json:"quantity"`
}
