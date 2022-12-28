package entities

type Race struct {
	Key            string          `json:"key"`
	Name           string          `json:"name"`
	Speed          int             `json:"speed"`
	AbilityBonuses []*AbilityBonus `json:"ability_bonuses"`
}

type AbilityScore struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type AbilityBonus struct {
	AbilityScore *AbilityScore `json:"ability_score"`
	Bonus        int           `json:"bonus"`
}
