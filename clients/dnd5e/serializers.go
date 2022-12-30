package dnd5e

import "github.com/fadedpez/dnd5e-api/entities"

func mapToReferenceItem(input map[string]interface{}) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Name: input["name"].(string),
		URL:  input["url"].(string),
		Key:  input["index"].(string),
	}
}

func mapToOption(input map[string]interface{}) entities.Option {
	if input == nil {
		return nil
	}

	switch input["option_type"].(string) {
	case "reference":
		return &entities.ReferenceOption{
			Reference: mapToReferenceItem(input["item"].(map[string]interface{})),
		}
	}

	return nil
}

func mapsToOptions(input []map[string]interface{}) []entities.Option {
	out := make([]entities.Option, len(input))
	for i, o := range input {
		out[i] = mapToOption(o)
	}

	return out
}

func mapToOptionSet(input map[string]interface{}) entities.OptionSet {
	if input == nil {
		return nil
	}

	switch input["option_set_type"].(string) {
	case "options_array":
		sliceMaps := make([]map[string]interface{}, len(input["options"].([]interface{})))
		sliceInterfaces := input["options"].([]interface{})
		for i, v := range sliceInterfaces {
			sliceMaps[i] = v.(map[string]interface{})
		}
		return &entities.OptionsArrayOptionSet{
			Options: mapsToOptions(sliceMaps),
		}
	}

	return nil
}

func listResultToRace(input *listResult) *entities.Race {
	return &entities.Race{
		Key:  input.Index,
		Name: input.Name,
	}
}

func abilityBonusResultToAbilityBonus(input *abilityBonus) *entities.AbilityBonus {
	if input == nil {
		return nil
	}

	return &entities.AbilityBonus{
		AbilityScore: &entities.AbilityScore{
			Key:  input.AbilityScore.Index,
			Name: input.AbilityScore.Name,
		},
		Bonus: input.Bonus,
	}
}

func abilityBonusResultsToAbilityBonuses(input []*abilityBonus) []*entities.AbilityBonus {
	out := make([]*entities.AbilityBonus, len(input))
	for i, b := range input {
		out[i] = abilityBonusResultToAbilityBonus(b)
	}

	return out
}

func languageResultToLanguage(input *language) *entities.Language {
	if input == nil {
		return nil
	}

	return &entities.Language{
		Key:  input.Index,
		Name: input.Name,
	}
}

func languageResultsToLanguages(input []*language) []*entities.Language {
	out := make([]*entities.Language, len(input))
	for i, l := range input {
		out[i] = languageResultToLanguage(l)
	}

	return out
}

func traitResultToTrait(input *trait) *entities.Trait {
	if input == nil {
		return nil
	}

	return &entities.Trait{
		Key:  input.Index,
		Name: input.Name,
	}
}

func traitResultsToTraits(input []*trait) []*entities.Trait {
	out := make([]*entities.Trait, len(input))
	for i, t := range input {
		out[i] = traitResultToTrait(t)
	}

	return out
}

func subRaceResultToSubRace(input *subRace) *entities.SubRace {
	if input == nil {
		return nil
	}

	return &entities.SubRace{
		Key:  input.Index,
		Name: input.Name,
	}
}

func subRaceResultsToSubRaces(input []*subRace) []*entities.SubRace {
	out := make([]*entities.SubRace, len(input))
	for i, s := range input {
		out[i] = subRaceResultToSubRace(s)
	}

	return out
}

func proficiencyResultToProficiency(input *proficiency) *entities.Proficiency {
	if input == nil {
		return nil
	}

	return &entities.Proficiency{
		Key:  input.Index,
		Name: input.Name,
	}
}

func proficiencyResultsToProficiencies(input []*proficiency) []*entities.Proficiency {
	out := make([]*entities.Proficiency, len(input))
	for i, p := range input {
		out[i] = proficiencyResultToProficiency(p)
	}

	return out
}

func listResultToEquipment(input *listResult) *entities.Equipment {
	return &entities.Equipment{
		Key:  input.Index,
		Name: input.Name,
	}
}

func equipmentCategoryResultToEquipmentCategory(input *equipmentCategory) *entities.EquipmentCategory {
	if input == nil {
		return nil
	}

	return &entities.EquipmentCategory{
		Key:  input.Index,
		Name: input.Name,
	}
}

func costResultToCost(input *cost) *entities.Cost {
	if input == nil {
		return nil
	}

	return &entities.Cost{
		Quantity: input.Quantity,
		Unit:     input.Unit,
	}
}

func equipmentResultToEquipment(input *equipmentResult) *entities.Equipment {
	if input == nil {
		return nil
	}

	return &entities.Equipment{
		Key:    input.Index,
		Name:   input.Name,
		Cost:   costResultToCost(input.Cost),
		Weight: input.Weight,
		EquipmentCategory: equipmentCategoryResultToEquipmentCategory(
			input.EquipmentCategory,
		),
	}
}

func damageResultToDamage(input *damage) *entities.Damage {
	if input == nil {
		return nil
	}

	return &entities.Damage{
		DamageDice: input.DamageDice,
		DamageType: damageTypeResultToDamageType(input.DamageType),
	}
}

func damageTypeResultToDamageType(input *damageType) *entities.DamageType {
	if input == nil {
		return nil
	}

	return &entities.DamageType{
		Key:  input.Index,
		Name: input.Name,
	}
}

func propertyResultToProperties(input *properties) *entities.Properties {
	if input == nil {
		return nil
	}

	return &entities.Properties{
		Key:  input.Index,
		Name: input.Name,
	}
}

func propertiesResultsToProperties(input []*properties) []*entities.Properties {
	out := make([]*entities.Properties, len(input))
	for i, p := range input {
		out[i] = propertyResultToProperties(p)
	}

	return out
}

func weaponRangeResultToWeaponRange(input *weaponRange) *entities.Range {
	if input == nil {
		return nil
	}

	return &entities.Range{
		Normal: input.Normal,
	}
}

func weaponResultToWeapon(input *weaponResult) *entities.Weapon {
	if input == nil {
		return nil
	}

	return &entities.Weapon{
		Key:               input.Index,
		Name:              input.Name,
		Cost:              costResultToCost(input.Cost),
		Damage:            damageResultToDamage(input.Damage),
		Weight:            input.Weight,
		EquipmentCategory: equipmentCategoryResultToEquipmentCategory(input.EquipmentCategory),
		WeaponCategory:    input.WeaponCategory,
		WeaponRange:       input.WeaponRange,
		CategoryRange:     input.CategoryRange,
		Range:             weaponRangeResultToWeaponRange(input.Range),
		Properties:        propertiesResultsToProperties(input.Properties),
		TwoHandedDamage:   damageResultToDamage(input.TwoHandedDamage),
	}
}

func armorResultToArmor(input *armorResult) *entities.Armor {
	if input == nil {
		return nil
	}

	return &entities.Armor{
		Key:                 input.Index,
		Name:                input.Name,
		Cost:                costResultToCost(input.Cost),
		ArmorCategory:       input.ArmorCategory,
		ArmorClass:          armorClassResultToArmorClass(input.ArmorClass),
		StrMinimum:          input.StrMinimum,
		StealthDisadvantage: input.StealthDisadvantage,
		Weight:              input.Weight,
		EquipmentCategory:   equipmentCategoryResultToEquipmentCategory(input.EquipmentCategory),
	}
}

func armorClassResultToArmorClass(input *armorClass) *entities.ArmorClass {
	if input == nil {
		return nil
	}

	return &entities.ArmorClass{
		Base:     input.Base,
		DexBonus: input.DexBonus,
	}
}

func listClassResultToClass(input *listResult) *entities.Class {
	return &entities.Class{
		Key:  input.Index,
		Name: input.Name,
	}
}

func savingThrowResultToSavingThrow(input *savingThrow) *entities.SavingThrow {
	if input == nil {
		return nil
	}

	return &entities.SavingThrow{
		Key:  input.Index,
		Name: input.Name,
	}
}

func savingThrowResultsToSavingThrows(input []*savingThrow) []*entities.SavingThrow {
	out := make([]*entities.SavingThrow, len(input))
	for i, s := range input {
		out[i] = savingThrowResultToSavingThrow(s)
	}

	return out
}

func equipmentListResultToEquipmentList(input *equipmentList) *entities.EquipmentList {
	if input == nil {
		return nil
	}

	return &entities.EquipmentList{
		Key:  input.Index,
		Name: input.Name,
	}
}

func startingEquipmentResultToStartingEquipment(input *startingEquipment) *entities.StartingEquipment {
	if input == nil {
		return nil
	}

	return &entities.StartingEquipment{
		Equipment: equipmentListResultToEquipmentList(input.Equipment),
		Quantity:  input.Quantity,
	}
}

func startingEquipmentResultsToStartingEquipment(input []*startingEquipment) []*entities.StartingEquipment {
	out := make([]*entities.StartingEquipment, len(input))
	for i, s := range input {
		out[i] = startingEquipmentResultToStartingEquipment(s)
	}

	return out
}
