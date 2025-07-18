package dnd5e

type referenceItem struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type raceResult struct {
	Index                      string           `json:"index"`
	Name                       string           `json:"name"`
	Speed                      int              `json:"speed"`
	Size                       string           `json:"size"`
	SizeDescription            string           `json:"size_description"`
	AbilityBonus               []*abilityBonus  `json:"ability_bonuses"`
	Language                   []*referenceItem `json:"languages"`
	Trait                      []*referenceItem `json:"traits"`
	SubRaces                   []*referenceItem `json:"subraces"`
	StartingProficiencies      []*referenceItem `json:"starting_proficiencies"`
	StartingProficiencyOptions *choiceResult    `json:"starting_proficiency_options"`
	LanguageOptions            *choiceResult    `json:"language_options"`
}

type abilityBonus struct {
	AbilityScore *referenceItem `json:"ability_score"`
	Bonus        int            `json:"bonus"`
}

type listResponse struct {
	Count   int              `json:"count"`
	Results []*referenceItem `json:"results"`
}

type equipmentListResponse struct {
	Equipment []*referenceItem `json:"equipment"`
}

type equipmentResult struct {
	Index             string         `json:"index"`
	Name              string         `json:"name"`
	Cost              *cost          `json:"cost"`
	Weight            float32        `json:"weight"`
	EquipmentCategory *referenceItem `json:"equipment_category"`
}

func (e *equipmentResult) getCategoryKey() string {
	if e.EquipmentCategory == nil {
		return ""
	}

	return e.EquipmentCategory.Index
}

type cost struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

type weaponResult struct {
	Index             string           `json:"index"`
	Name              string           `json:"name"`
	Cost              *cost            `json:"cost"`
	Weight            float32          `json:"weight"`
	EquipmentCategory *referenceItem   `json:"equipment_category"`
	WeaponCategory    string           `json:"weapon_category"`
	WeaponRange       string           `json:"weapon_range"`
	CategoryRange     string           `json:"category_range"`
	Damage            *damage          `json:"damage"`
	Range             *weaponRange     `json:"range"`
	Properties        []*referenceItem `json:"properties"`
	TwoHandedDamage   *damage          `json:"two_handed_damage"`
}

type damage struct {
	DamageDice string         `json:"damage_dice"`
	DamageType *referenceItem `json:"damage_type"`
}

type weaponRange struct {
	Normal int `json:"normal"`
}

type armorResult struct {
	Index               string         `json:"index"`
	Name                string         `json:"name"`
	Cost                *cost          `json:"cost"`
	Weight              float32        `json:"weight"`
	EquipmentCategory   *referenceItem `json:"equipment_category"`
	ArmorCategory       string         `json:"armor_category"`
	ArmorClass          *armorClass    `json:"armor_class"`
	StrMinimum          int            `json:"str_minimum"`
	StealthDisadvantage bool           `json:"stealth_disadvantage"`
}

type armorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
}

type classResult struct {
	Index                    string               `json:"index"`
	Name                     string               `json:"name"`
	HitDie                   int                  `json:"hit_die"`
	Proficiencies            []*referenceItem     `json:"proficiencies"`
	SavingThrows             []*referenceItem     `json:"saving_throws"`
	StartingEquipment        []*startingEquipment `json:"starting_equipment"`
	ProficiencyChoices       []*choiceResult      `json:"proficiency_choices"`
	StartingEquipmentOptions []*choiceResult      `json:"starting_equipment_options"`
	MultiClassing            *multiClassing       `json:"multi_classing"`
}

type multiClassing struct {
	Prerequisites       []*multiClassingPrerequisite `json:"prerequisites"`
	Proficiencies       []*referenceItem             `json:"proficiencies"`
	ProficiencyChoices  []*choiceResult              `json:"proficiency_choices"`
}

type multiClassingPrerequisite struct {
	AbilityScore *referenceItem `json:"ability_score"`
	MinimumScore int            `json:"minimum_score"`
}

type startingEquipment struct {
	Equipment *referenceItem `json:"equipment"`
	Quantity  int            `json:"quantity"`
}

type spellResult struct {
	Index         string           `json:"index"`
	Name          string           `json:"name"`
	Range         string           `json:"range"`
	Ritual        bool             `json:"ritual"`
	Duration      string           `json:"duration"`
	Concentration bool             `json:"concentration"`
	CastingTime   string           `json:"casting_time"`
	SpellLevel    int              `json:"level"`
	SpellDamage   *spellDamage     `json:"damage"`
	DC            *dc              `json:"dc"`
	AreaOfEffect  *areaOfEffect    `json:"area_of_effect"`
	SpellSchool   *referenceItem   `json:"school"`
	SpellClasses  []*referenceItem `json:"classes"`
}

type spellDamage struct {
	DamageType        *referenceItem          `json:"damage_type"`
	DamageAtSlotLevel *spellDamageAtSlotLevel `json:"damage_at_slot_level"`
}

type spellDamageAtSlotLevel struct {
	FirstLevel   string `json:"1"`
	SecondLevel  string `json:"2"`
	ThirdLevel   string `json:"3"`
	FourthLevel  string `json:"4"`
	FifthLevel   string `json:"5"`
	SixthLevel   string `json:"6"`
	SeventhLevel string `json:"7"`
	EighthLevel  string `json:"8"`
	NinthLevel   string `json:"9"`
}

type dc struct {
	DCType    *referenceItem `json:"dc_type"`
	DCSuccess string         `json:"dc_success"`
}

type areaOfEffect struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}

type featureResult struct {
	Index           string            `json:"index"`
	Name            string            `json:"name"`
	Level           int               `json:"level"`
	Class           *referenceItem    `json:"class"`
	FeatureSpecific *subFeatureOption `json:"feature_specific"`
	Invocations     []*referenceItem  `json:"invocations"`
}

type subFeatureOption struct {
	SubFeatureOptions *choiceResult `json:"subfeature_options"`
}

type skillResult struct {
	Index        string         `json:"index"`
	Name         string         `json:"name"`
	Description  []string       `json:"desc"`
	AbilityScore *referenceItem `json:"ability_score"`
	URL          string         `json:"url"`
}

type monsterResult struct {
	Index                 string                `json:"index"`
	Name                  string                `json:"name"`
	Size                  string                `json:"size"`
	Type                  string                `json:"type"`
	Alignment             string                `json:"alignment"`
	ArmorClass            []*monsterArmorClass  `json:"armor_class"`
	HitPoints             int                   `json:"hit_points"`
	HitDice               string                `json:"hit_dice"`
	HitPointsRoll         string                `json:"hit_points_roll"`
	Speed                 *monsterSpeed         `json:"speed"`
	Strength              int                   `json:"strength"`
	Dexterity             int                   `json:"dexterity"`
	Constitution          int                   `json:"constitution"`
	Intelligence          int                   `json:"intelligence"`
	Wisdom                int                   `json:"wisdom"`
	Charisma              int                   `json:"charisma"`
	Proficiencies         []*monsterProficiency `json:"proficiencies"`
	DamageVulnerabilities []string              `json:"damage_vulnerabilities"`
	DamageResistances     []string              `json:"damage_resistances"`
	DamageImmunities      []string              `json:"damage_immunities"`
	ConditionImmunities   []*referenceItem      `json:"condition_immunities"`
	Senses                *monsterSenses        `json:"senses"`
	Languages             string                `json:"languages"`
	ChallengeRating       float32               `json:"challenge_rating"`
	XP                    int                   `json:"xp"`
	MonsterActions        []*monsterAction      `json:"actions"` //TODO: convert to an interface
	MonsterImageURL       string                `json:"image"`
	//TODO: Add legendary actions
	//TODO: Add reactions
	//TODO: Add special abilities, possible to be an interface?
}

type monsterArmorClass struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type monsterSpeed struct {
	Walk   string `json:"walk"`
	Burrow string `json:"burrow"`
	Climb  string `json:"climb"`
	Fly    string `json:"fly"`
	Swim   string `json:"swim"`
}

type monsterProficiency struct {
	Value       int            `json:"value"`
	Proficiency *referenceItem `json:"proficiency"`
}

type monsterSenses struct {
	Blindsight        string `json:"blindsight"`
	Darkvision        string `json:"darkvision"`
	Tremorsense       string `json:"tremorsense"`
	Truesight         string `json:"truesight"`
	PassivePerception int    `json:"passive_perception"`
}

type monsterAction struct {
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	AttackBonus int       `json:"attack_bonus"`
	Damage      []*damage `json:"damage"`
}

type levelResult struct {
	Level               int                  `json:"level"`
	AbilityScoreBonuses int                  `json:"ability_score_bonuses"`
	ProfBonus           int                  `json:"prof_bonus"`
	Features            []*referenceItem     `json:"features"`
	SpellCasting        *spellCasting        `json:"spellcasting"`
	ClassSpecific       *classSpecificResult `json:"class_specific"`
	Index               string               `json:"index"`
	Class               *referenceItem       `json:"class"`
}

type spellCasting struct {
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

type classSpecificResult struct {
	FavoredEnemies         int                   `json:"favored_enemies"`
	FavoredTerrain         int                   `json:"favored_terrain"`
	RageCount              int                   `json:"rage_count"`
	RageDamageBonus        int                   `json:"rage_damage_bonus"`
	BrutalCriticalDice     int                   `json:"brutal_critical_dice"`
	BardicInspirationDie   int                   `json:"bardic_inspiration_die"`
	SongOfRestDie          int                   `json:"song_of_rest_die"`
	MagicalSecretsMax5     int                   `json:"magical_secrets_max_5"`
	MagicalSecretsMax7     int                   `json:"magical_secrets_max_7"`
	MagicalSecretsMax9     int                   `json:"magical_secrets_max_9"`
	ChannelDivinityCharges int                   `json:"channel_divinity_charges"`
	DestroyUndeadCR        int                   `json:"destroy_undead_cr"`
	WildShapeMaxCR         int                   `json:"wild_shape_max_cr"`
	WildShapeSwim          bool                  `json:"wild_shape_swim"`
	WildShapeFly           bool                  `json:"wild_shape_fly"`
	ActionSurges           int                   `json:"action_surges"`
	IndomitableUses        int                   `json:"indomitable_uses"`
	ExtraAttacks           int                   `json:"extra_attacks"`
	MartialArts            *martialArts          `json:"martial_arts"`
	KiPoints               int                   `json:"ki_points"`
	UnarmoredMovement      int                   `json:"unarmored_movement"`
	AuraRange              int                   `json:"aura_range"`
	SneakAttack            *sneakAttack          `json:"sneak_attack"`
	SorceryPoints          int                   `json:"sorcery_points"`
	MetamagicKnown         int                   `json:"metamagic_known"`
	CreatingSpellSlots     []*creatingSpellSlots `json:"creating_spell_slots"`
	InvocationsKnown       int                   `json:"invocations_known"`
	MysticArcanumLevel6    int                   `json:"mystic_arcanum_level_6"`
	MysticArcanumLevel7    int                   `json:"mystic_arcanum_level_7"`
	MysticArcanumLevel8    int                   `json:"mystic_arcanum_level_8"`
	MysticArcanumLevel9    int                   `json:"mystic_arcanum_level_9"`
	ArcaneRecoveryLevels   int                   `json:"arcane_recovery_levels"`
}

type martialArts struct {
	DiceCount int `json:"dice_count"`
	DiceValue int `json:"dice_value"`
}

type sneakAttack struct {
	DiceCount int `json:"dice_count"`
	DiceValue int `json:"dice_value"`
}

type creatingSpellSlots struct {
	SpellSlotLevel   int `json:"spell_slot_level"`
	SorceryPointCost int `json:"sorcery_point_cost"`
}

type proficiencyResult struct {
	Index     string         `json:"index"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Reference *referenceItem `json:"reference"`
}

type damageTypeResult struct {
	Index       string   `json:"index"`
	Name        string   `json:"name"`
	Description []string `json:"desc"`
	URL         string   `json:"url"`
}

type backgroundResult struct {
	Index                      string           `json:"index"`
	Name                       string           `json:"name"`
	StartingProficiencies      []*referenceItem `json:"starting_proficiencies"`
	LanguageOptions            *choiceResult    `json:"language_options"`
	StartingEquipment          []*startingEquipment `json:"starting_equipment"`
	StartingEquipmentOptions   []*choiceResult  `json:"starting_equipment_options"`
	Feature                    *backgroundFeatureResult `json:"feature"`
	PersonalityTraits          *choiceResult    `json:"personality_traits"`
	Ideals                     *choiceResult    `json:"ideals"`
	Bonds                      *choiceResult    `json:"bonds"`
	Flaws                      *choiceResult    `json:"flaws"`
}

type backgroundFeatureResult struct {
	Name        string   `json:"name"`
	Description []string `json:"desc"`
}
