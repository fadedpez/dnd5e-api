package entities

type Race struct {
	Key            string          `json:"key"`
	Name           string          `json:"name"`
	Speed          int             `json:"speed"`
	AbilityBonuses []*AbilityBonus `json:"ability_bonuses"`
	Languages      []*Language     `json:"languages"`
	Traits         []*Trait        `json:"traits"`
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
