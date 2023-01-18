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
	case "counted_reference":
		return &entities.ReferenceOption{
			Reference: mapToReferenceItem(input["of"].(map[string]interface{})),
		}
	case "choice":
		return mapToChoice(input["choice"].(map[string]interface{}))
	}

	return nil
}

func mapsToChoices(input []map[string]interface{}) []*entities.Choice {
	out := make([]*entities.Choice, len(input))
	for i, c := range input {
		out[i] = mapToChoice(c)
	}

	return out
}

func mapToChoice(input map[string]interface{}) *entities.Choice {
	if input == nil {
		return nil
	}

	return &entities.Choice{
		Choose:    int(input["choose"].(float64)),
		Type:      input["type"].(string),
		OptionSet: mapToOptionList(input["from"].(map[string]interface{})),
	}
}

func mapsToOptions(input []map[string]interface{}) []entities.Option {
	out := make([]entities.Option, len(input))
	for i, o := range input {
		out[i] = mapToOption(o)
	}

	return out
}

func mapToOptionList(input map[string]interface{}) *entities.OptionList {
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
		return &entities.OptionList{
			Options: mapsToOptions(sliceMaps),
		}
	case "equipment_category":

	}

	return nil
}

func listResultToRace(input *listResult) *entities.ReferenceItem {
	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func abilityBonusResultToAbilityBonus(input *abilityBonus) *entities.AbilityBonus {
	if input == nil {
		return nil
	}

	return &entities.AbilityBonus{
		AbilityScore: &entities.ReferenceItem{
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

func languageResultToLanguage(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func languageResultsToLanguages(input []*listResult) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, l := range input {
		out[i] = languageResultToLanguage(l)
	}

	return out
}

func traitResultToTrait(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func traitResultsToTraits(input []*listResult) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, t := range input {
		out[i] = traitResultToTrait(t)
	}

	return out
}

func subRaceResultToSubRace(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func subRaceResultsToSubRaces(input []*listResult) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, s := range input {
		out[i] = subRaceResultToSubRace(s)
	}

	return out
}

func proficiencyResultToProficiency(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func proficiencyResultsToProficiencies(input []*listResult) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, p := range input {
		out[i] = proficiencyResultToProficiency(p)
	}

	return out
}

func listResultToEquipment(input *listResult) *entities.ReferenceItem {
	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func equipmentCategoryResultToEquipmentCategory(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
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

func damageTypeResultToDamageType(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func propertyResultToProperties(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func propertiesResultsToProperties(input []*listResult) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
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

func listClassResultToClass(input *listResult) *entities.ReferenceItem {
	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func savingThrowResultToSavingThrow(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func savingThrowResultsToSavingThrows(input []*listResult) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, s := range input {
		out[i] = savingThrowResultToSavingThrow(s)
	}

	return out
}

func equipmentListResultToEquipmentList(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
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

func listResultToSpell(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func spellDamageResultToSpellDamage(input *spellDamage) *entities.SpellDamage {
	if input == nil {
		return nil
	}

	return &entities.SpellDamage{
		SpellDamageType:        spellDamageTypeResultToSpellDamageType(input.DamageType),
		SpellDamageAtSlotLevel: spellDamageAtSlotLevelToSpellDamageAtSlotLevel(input.DamageAtSlotLevel),
	}
}

func spellDamageTypeResultToSpellDamageType(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func spellDamageAtSlotLevelToSpellDamageAtSlotLevel(input *spellDamageAtSlotLevel) *entities.SpellDamageAtSlotLevel {
	if input == nil {
		return nil
	}

	return &entities.SpellDamageAtSlotLevel{
		FirstLevel:   input.FirstLevel,
		SecondLevel:  input.SecondLevel,
		ThirdLevel:   input.ThirdLevel,
		FourthLevel:  input.FourthLevel,
		FifthLevel:   input.FifthLevel,
		SixthLevel:   input.SixthLevel,
		SeventhLevel: input.SeventhLevel,
		EighthLevel:  input.EighthLevel,
		NinthLevel:   input.NinthLevel,
	}
}

func dcResultToDC(input *dc) *entities.DC {
	if input == nil {
		return nil
	}

	return &entities.DC{
		DCType:    dcTypeResultToDCType(input.DCType),
		DCSuccess: input.DCSuccess,
	}
}

func dcTypeResultToDCType(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func spellClassResultsToSpellClasses(input []*listResult) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, s := range input {
		out[i] = &entities.ReferenceItem{
			Key:  s.Index,
			Name: s.Name,
		}
	}

	return out
}

func areaOfEffectResultToAreaOfEffect(input *areaOfEffect) *entities.AreaOfEffect {
	if input == nil {
		return nil
	}

	return &entities.AreaOfEffect{
		Type: input.Type,
		Size: input.Size,
	}
}

func spellSchoolResultToSpellSchool(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func listResultToFeature(input *listResult) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}
