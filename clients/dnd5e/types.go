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
	Weight            int            `json:"weight"`
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
	Weight            int              `json:"weight"`
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
	Weight              int            `json:"weight"`
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
