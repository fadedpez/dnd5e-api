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
	CantripsKnown    int `json:"cantrips_known"`
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

type BardSpecific struct {
	BardicInspirationDie int `json:"bardic_inspiration_die"`
	SongOfRestDie        int `json:"song_of_rest_die"`
	MagicalSecretsMax5   int `json:"magical_secret_max_5"`
	MagicalSecretsMax7   int `json:"magical_secret_max_7"`
	MagicalSecretsMax9   int `json:"magic_secret_max_9"`
}

func (b BardSpecific) GetSpecificClass() string {
	return "bard"
}

type ClericSpecific struct {
	ChannelDivinityCharges int `json:"channel_divinity_charges"`
	DestroyUndeadCR        int `json:"destroy_undead_cr"`
}

func (c ClericSpecific) GetSpecificClass() string {
	return "cleric"
}

type DruidSpecific struct {
	WildShapeMaxCR int  `json:"wild_shape_max_cr"`
	WildShapeSwim  bool `json:"wild_shape_swim"`
	WildShapeFly   bool `json:"wild_shape_fly"`
}

func (d DruidSpecific) GetSpecificClass() string {
	return "druid"
}

type FighterSpecific struct {
	ActionSurges    int `json:"action_surges"`
	IndomitableUses int `json:"indomitable_uses"`
	ExtraAttacks    int `json:"extra_attacks"`
}

func (f FighterSpecific) GetSpecificClass() string {
	return "fighter"
}

type MonkSpecific struct {
	MartialArts       *MartialArts `json:"martial_arts"`
	KiPoints          int          `json:"ki_points"`
	UnarmoredMovement int          `json:"unarmored_movement"`
}

func (m MonkSpecific) GetSpecificClass() string {
	return "monk"
}

type MartialArts struct {
	DiceCount int `json:"dice_count"`
	DiceValue int `json:"dice_value"`
}

type PaladinSpecific struct {
	AuraRange int `json:"aura_range"`
}

func (p PaladinSpecific) GetSpecificClass() string {
	return "paladin"
}

type RogueSpecific struct {
	SneakAttack *SneakAttack `json:"sneak_attack_dice"`
}

func (r RogueSpecific) GetSpecificClass() string {
	return "rogue"
}

type SneakAttack struct {
	DiceCount int `json:"dice_count"`
	DiceValue int `json:"dice_value"`
}

type SorcererSpecific struct {
	SorceryPoints      int                   `json:"sorcery_points"`
	MetamagicKnown     int                   `json:"metamagic_known"`
	CreatingSpellSlots []*CreatingSpellSlots `json:"creating_spell_slots"`
}

func (s SorcererSpecific) GetSpecificClass() string {
	return "sorcerer"
}

type CreatingSpellSlots struct {
	SpellSlotLevel   int `json:"spell_slot_level"`
	SorceryPointCost int `json:"sorcery_point_cost"`
}

type WarlockSpecific struct {
	InvocationsKnown    int `json:"invocations_known"`
	MysticArcanumLevel6 int `json:"mystic_arcanum_level_6"`
	MysticArcanumLevel7 int `json:"mystic_arcanum_level_7"`
	MysticArcanumLevel8 int `json:"mystic_arcanum_level_8"`
	MysticArcanumLevel9 int `json:"mystic_arcanum_level_9"`
}

func (w WarlockSpecific) GetSpecificClass() string {
	return "warlock"
}

type WizardSpecific struct {
	ArcaneRecoveryLevels int `json:"arcane_recovery_levels"`
}

func (w WizardSpecific) GetSpecificClass() string {
	return "wizard"
}
