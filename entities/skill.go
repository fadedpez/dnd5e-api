package entities

type Skill struct {
	Key          string         `json:"index"`
	Name         string         `json:"name"`
	Description []string       `json:"desc"`
	AbilityScore *ReferenceItem `json:"ability_score"`
	Type         string         `json:"type"`
}
