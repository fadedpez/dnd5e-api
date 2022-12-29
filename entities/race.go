package entities

type Race struct {
	Key                        string          `json:"key"`
	Name                       string          `json:"name"`
	Speed                      int             `json:"speed"`
	AbilityBonuses             []*AbilityBonus `json:"ability_bonuses"`
	Languages                  []*Language     `json:"languages"`
	Traits                     []*Trait        `json:"traits"`
	SubRaces                   []*SubRace      `json:"subrace"`
	StartingProficiencies      []*Proficiency  `json:"starting_proficiencies"`
	StartingProficiencyOptions *Choice         `json:"starting_proficiency_options"`
	LanguageOptions            *Choice         `json:"language_options"`
}

type AbilityScore struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type AbilityBonus struct {
	AbilityScore *AbilityScore `json:"ability_score"`
	Bonus        int           `json:"bonus"`
}

type Language struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Trait struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type SubRace struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Proficiency struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
