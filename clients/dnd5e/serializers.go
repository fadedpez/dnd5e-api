package dnd5e

import (
	"github.com/fadedpez/dnd5e-api/entities"
)

func referenceItemToRace(input *referenceItem) *entities.ReferenceItem {
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

func languageResultToLanguage(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func languageResultsToLanguages(input []*referenceItem) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, l := range input {
		out[i] = languageResultToLanguage(l)
	}

	return out
}

func traitResultToTrait(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func traitResultsToTraits(input []*referenceItem) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, t := range input {
		out[i] = traitResultToTrait(t)
	}

	return out
}

func subRaceResultToSubRace(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func subRaceResultsToSubRaces(input []*referenceItem) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, s := range input {
		out[i] = subRaceResultToSubRace(s)
	}

	return out
}

func proficiencyResultToProficiency(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func proficiencyResultsToProficiencies(input []*referenceItem) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, p := range input {
		out[i] = proficiencyResultToProficiency(p)
	}

	return out
}

func referenceItemToEquipment(input *referenceItem) *entities.ReferenceItem {
	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func equipmentCategoryResultToEquipmentCategory(input *referenceItem) *entities.ReferenceItem {
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

func damageTypeResultToDamageType(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func propertyResultToProperties(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func propertiesResultsToProperties(input []*referenceItem) []*entities.ReferenceItem {
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

func referenceItemToClass(input *referenceItem) *entities.ReferenceItem {
	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func savingThrowResultToSavingThrow(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func savingThrowResultsToSavingThrows(input []*referenceItem) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, s := range input {
		out[i] = savingThrowResultToSavingThrow(s)
	}

	return out
}

func equipmentListResultToEquipmentList(input *referenceItem) *entities.ReferenceItem {
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

func referenceItemToSpell(input *referenceItem) *entities.ReferenceItem {
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

func spellDamageTypeResultToSpellDamageType(input *referenceItem) *entities.ReferenceItem {
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

func dcTypeResultToDCType(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func spellClassResultsToSpellClasses(input []*referenceItem) []*entities.ReferenceItem {
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

func spellSchoolResultToSpellSchool(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func referenceItemToFeature(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func featureClassResultToClass(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func choiceResultToChoice(input *choiceResult) *entities.ChoiceOption {
	if input == nil {
		return nil
	}

	return input.toEntity()
}

func choiceResultsToChoices(input []*choiceResult) []*entities.ChoiceOption {
	if input == nil {
		return nil
	}

	out := make([]*entities.ChoiceOption, len(input))
	for i, c := range input {
		out[i] = choiceResultToChoice(c)
	}

	return out
}

func referenceItemToSkill(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}
