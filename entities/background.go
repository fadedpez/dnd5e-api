package entities

type Background struct {
	Key                string                   `json:"key"`
	Name               string                   `json:"name"`
	SkillProficiencies []*ReferenceItem         `json:"skill_proficiencies"`
	LanguageOptions    *ChoiceOption            `json:"language_options"`
	StartingEquipment  []*StartingEquipment     `json:"starting_equipment"`
	StartingEquipmentOptions []*ChoiceOption  `json:"starting_equipment_options"`
	Feature            *BackgroundFeature       `json:"feature"`
	PersonalityTraits  *ChoiceOption            `json:"personality_traits"`
	Ideals             *ChoiceOption            `json:"ideals"`
	Bonds              *ChoiceOption            `json:"bonds"`
	Flaws              *ChoiceOption            `json:"flaws"`
}

type BackgroundFeature struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}