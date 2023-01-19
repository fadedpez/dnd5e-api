package entities

type Race struct {
	Key                        string           `json:"key"`
	Name                       string           `json:"name"`
	Speed                      int              `json:"speed"`
	AbilityBonuses             []*AbilityBonus  `json:"ability_bonuses"`
	Languages                  []*ReferenceItem `json:"languages"`
	Traits                     []*ReferenceItem `json:"traits"`
	SubRaces                   []*ReferenceItem `json:"subrace"`
	StartingProficiencies      []*ReferenceItem `json:"starting_proficiencies"`
	StartingProficiencyOptions *ChoiceOption    `json:"starting_proficiency_options"`
	LanguageOptions            *ChoiceOption    `json:"language_options"`
}

type AbilityBonus struct {
	AbilityScore *ReferenceItem `json:"ability_score"`
	Bonus        int            `json:"bonus"`
}
