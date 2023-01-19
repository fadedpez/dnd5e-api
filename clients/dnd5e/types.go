package dnd5e

import "github.com/fadedpez/dnd5e-api/entities"

type listResult struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type raceResult struct {
	Index                      string                 `json:"index"`
	Name                       string                 `json:"name"`
	Speed                      int                    `json:"speed"`
	AbilityBonus               []*abilityBonus        `json:"ability_bonuses"`
	Language                   []*listResult          `json:"languages"`
	Trait                      []*listResult          `json:"traits"`
	SubRaces                   []*listResult          `json:"subraces"`
	StartingProficiencies      []*listResult          `json:"starting_proficiencies"`
	StartingProficiencyOptions map[string]interface{} `json:"starting_proficiency_options"`
	LanguageOptions            map[string]interface{} `json:"language_options"`
}

func (r *raceResult) getStartingProficiencyChoice() *entities.Choice {
	if r.StartingProficiencyOptions == nil {
		return nil
	}

	out := &entities.Choice{
		Choose:    int(r.StartingProficiencyOptions["choose"].(float64)),
		Type:      r.StartingProficiencyOptions["type"].(string),
		OptionSet: mapToOptionList(r.StartingProficiencyOptions["from"].(map[string]interface{})),
	}

	return out
}

func (r *raceResult) getLanguageChoice() *entities.Choice {
	if r.LanguageOptions == nil {
		return nil
	}

	out := &entities.Choice{
		Choose:    int(r.LanguageOptions["choose"].(float64)),
		Type:      r.LanguageOptions["type"].(string),
		OptionSet: mapToOptionList(r.LanguageOptions["from"].(map[string]interface{})),
	}

	return out
}

type abilityBonus struct {
	AbilityScore *listResult `json:"ability_score"`
	Bonus        int         `json:"bonus"`
}

type listResponse struct {
	Count   int           `json:"count"`
	Results []*listResult `json:"results"`
}

type equipmentListResponse struct {
	Equipment []*listResult `json:"equipment"`
}

type equipmentResult struct {
	Index             string      `json:"index"`
	Name              string      `json:"name"`
	Cost              *cost       `json:"cost"`
	Weight            int         `json:"weight"`
	EquipmentCategory *listResult `json:"equipment_category"`
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
	Index             string        `json:"index"`
	Name              string        `json:"name"`
	Cost              *cost         `json:"cost"`
	Weight            int           `json:"weight"`
	EquipmentCategory *listResult   `json:"equipment_category"`
	WeaponCategory    string        `json:"weapon_category"`
	WeaponRange       string        `json:"weapon_range"`
	CategoryRange     string        `json:"category_range"`
	Damage            *damage       `json:"damage"`
	Range             *weaponRange  `json:"range"`
	Properties        []*listResult `json:"properties"`
	TwoHandedDamage   *damage       `json:"two_handed_damage"`
}

type damage struct {
	DamageDice string      `json:"damage_dice"`
	DamageType *listResult `json:"damage_type"`
}

type weaponRange struct {
	Normal int `json:"normal"`
}

type armorResult struct {
	Index               string      `json:"index"`
	Name                string      `json:"name"`
	Cost                *cost       `json:"cost"`
	Weight              int         `json:"weight"`
	EquipmentCategory   *listResult `json:"equipment_category"`
	ArmorCategory       string      `json:"armor_category"`
	ArmorClass          *armorClass `json:"armor_class"`
	StrMinimum          int         `json:"str_minimum"`
	StealthDisadvantage bool        `json:"stealth_disadvantage"`
}

type armorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
}

type classResult struct {
	Index                    string               `json:"index"`
	Name                     string               `json:"name"`
	HitDie                   int                  `json:"hit_die"`
	Proficiencies            []*listResult        `json:"proficiencies"`
	SavingThrows             []*listResult        `json:"saving_throws"`
	StartingEquipment        []*startingEquipment `json:"starting_equipment"`
	ProficiencyChoices       []*choiceResult      `json:"proficiency_choices"`
	StartingEquipmentOptions []*choiceResult      `json:"starting_equipment_options"`
}

type startingEquipment struct {
	Equipment *listResult `json:"equipment"`
	Quantity  int         `json:"quantity"`
}

type spellResult struct {
	Index         string        `json:"index"`
	Name          string        `json:"name"`
	Range         string        `json:"range"`
	Ritual        bool          `json:"ritual"`
	Duration      string        `json:"duration"`
	Concentration bool          `json:"concentration"`
	CastingTime   string        `json:"casting_time"`
	SpellLevel    int           `json:"level"`
	SpellDamage   *spellDamage  `json:"damage"`
	DC            *dc           `json:"dc"`
	AreaOfEffect  *areaOfEffect `json:"area_of_effect"`
	SpellSchool   *listResult   `json:"school"`
	SpellClasses  []*listResult `json:"classes"`
}

type spellDamage struct {
	DamageType        *listResult             `json:"damage_type"`
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
	DCType    *listResult `json:"dc_type"`
	DCSuccess string      `json:"dc_success"`
}

type areaOfEffect struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}

type featureResult struct {
	Index           string      `json:"index"`
	Name            string      `json:"name"`
	Level           int         `json:"level"`
	Class           *listResult `json:"class"`
	FeatureSpecific *subFeature `json:"feature_specific"`
}

type subFeature struct {
	SubfeatureOptions []map[string]interface{} `json:"subfeature_options"`
	Invocations       []*listResult            `json:"invocations"`
}
