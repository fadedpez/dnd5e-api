package entities

type Class struct {
	Key           string         `json:"key"`
	Name          string         `json:"name"`
	HitDie        int            `json:"hit_die"`
	Proficiencies []*Proficiency `json:"proficiencies"`
	SavingThrows  []*SavingThrow `json:"saving_throws"`
}

type SavingThrow struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
