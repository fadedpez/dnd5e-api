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
	Language                   []*language            `json:"languages"`
	Trait                      []*trait               `json:"traits"`
	SubRaces                   []*subRace             `json:"subraces"`
	StartingProficiencies      []*proficiency         `json:"starting_proficiencies"`
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
		OptionSet: mapToOptionSet(r.StartingProficiencyOptions["from"].(map[string]interface{})),
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
		OptionSet: mapToOptionSet(r.LanguageOptions["from"].(map[string]interface{})),
	}

	return out
}

type abilityBonus struct {
	AbilityScore *listResult `json:"ability_score"`
	Bonus        int         `json:"bonus"`
}

type language struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type trait struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type subRace struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type proficiency struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type listResponse struct {
	Count   int           `json:"count"`
	Results []*listResult `json:"results"`
}

type equipmentResult struct {
	Index             string             `json:"index"`
	Name              string             `json:"name"`
	Cost              *cost              `json:"cost"`
	Weight            int                `json:"weight"`
	EquipmentCategory *equipmentCategory `json:"equipment_category"`
}

func (e *equipmentResult) getCategoryKey() string {
	if e.EquipmentCategory == nil {
		return ""
	}

	return e.EquipmentCategory.Index
}

type equipmentCategory struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type cost struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}

type weaponResult struct {
	Index             string             `json:"index"`
	Name              string             `json:"name"`
	Cost              *cost              `json:"cost"`
	Weight            int                `json:"weight"`
	EquipmentCategory *equipmentCategory `json:"equipment_category"`
	WeaponCategory    string             `json:"weapon_category"`
	WeaponRange       string             `json:"weapon_range"`
	CategoryRange     string             `json:"category_range"`
	Damage            *damage            `json:"damage"`
	Range             *weaponRange       `json:"range"`
	Properties        []*properties      `json:"properties"`
	TwoHandedDamage   *damage            `json:"two_handed_damage"`
}

type damage struct {
	DamageDice string      `json:"damage_dice"`
	DamageType *damageType `json:"damage_type"`
}

type damageType struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type weaponRange struct {
	Normal int `json:"normal"`
}

type properties struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type armorResult struct {
	Index               string             `json:"index"`
	Name                string             `json:"name"`
	Cost                *cost              `json:"cost"`
	Weight              int                `json:"weight"`
	EquipmentCategory   *equipmentCategory `json:"equipment_category"`
	ArmorCategory       string             `json:"armor_category"`
	ArmorClass          *armorClass        `json:"armor_class"`
	StrMinimum          int                `json:"str_minimum"`
	StealthDisadvantage bool               `json:"stealth_disadvantage"`
}

type armorClass struct {
	Base     int  `json:"base"`
	DexBonus bool `json:"dex_bonus"`
}

type classResult struct {
	Index  string `json:"index"`
	Name   string `json:"name"`
	HitDie int    `json:"hit_die"`
}
