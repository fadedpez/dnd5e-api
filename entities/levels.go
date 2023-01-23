package entities

type Level struct {
	Level               int              `json:"level"`
	AbilityScoreBonuses int              `json:"ability_score_bonuses"`
	ProfBonus           int              `json:"prof_bonus"`
	Features            []*ReferenceItem `json:"features"`
	SpellCasting        *SpellCasting    `json:"spellcasting"`
	ClassSpecific       ClassSpecific    `json:"class_specific"`
	Key                 string           `json:"index"`
	Class               *ReferenceItem   `json:"class"`
}

type SpellCasting struct {
	SpellsKnown      int `json:"spells_known"`
	SpellSlotsLevel1 int `json:"spell_slots_level_1"`
	SpellSlotsLevel2 int `json:"spell_slots_level_2"`
	SpellSlotsLevel3 int `json:"spell_slots_level_3"`
	SpellSlotsLevel4 int `json:"spell_slots_level_4"`
	SpellSlotsLevel5 int `json:"spell_slots_level_5"`
	SpellSlotsLevel6 int `json:"spell_slots_level_6"`
	SpellSlotsLevel7 int `json:"spell_slots_level_7"`
	SpellSlotsLevel8 int `json:"spell_slots_level_8"`
	SpellSlotsLevel9 int `json:"spell_slots_level_9"`
}

type ClassSpecific interface {
	GetSpecificClass() string
}

type RangerSpecific struct {
	FavoredEnemies int `json:"favored_enemies"`
	FavoredTerrain int `json:"favored_terrain"`
}

func (r RangerSpecific) GetSpecificClass() string {
	return "ranger"
}

type BarbarianSpecific struct {
	RageCount          int `json:"rage_count"`
	RageDamageBonus    int `json:"rage_damage_bonus"`
	BrutalCriticalDice int `json:"brutal_critical_dice"`
}

func (b BarbarianSpecific) GetSpecificClass() string {
	return "barbarian"
}
