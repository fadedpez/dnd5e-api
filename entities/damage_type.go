package entities

type DamageType struct {
	Key         string   `json:"index"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Description []string `json:"desc"`
}
