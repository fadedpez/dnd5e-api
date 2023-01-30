package dnd5e

import (
	"github.com/fadedpez/dnd5e-api/entities"
)

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
		EquipmentCategory: referenceItemToReferenceItem(
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
		EquipmentCategory: referenceItemToReferenceItem(input.EquipmentCategory),
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
		EquipmentCategory:   referenceItemToReferenceItem(input.EquipmentCategory),
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

func referenceItemsToFeatures(input []*referenceItem) []*entities.ReferenceItem {
	out := make([]*entities.ReferenceItem, len(input))
	for i, s := range input {
		out[i] = referenceItemToFeature(s)
	}

	return out
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

func referenceItemToAbilityScore(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func referenceItemToMonster(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
	}
}

func monsterSpeedResultToSpeed(input *monsterSpeed) *entities.Speed {
	if input == nil {
		return nil
	}

	return &entities.Speed{
		Walk:   input.Walk,
		Climb:  input.Climb,
		Fly:    input.Fly,
		Swim:   input.Swim,
		Burrow: input.Burrow,
	}
}

func monsterProficiencyResultToMonsterProficiency(input *monsterProficiency) *entities.MonsterProficiency {
	if input == nil {
		return nil
	}

	return &entities.MonsterProficiency{
		Proficiency: referenceItemToReferenceItem(input.Proficiency),
		Value:       input.Value,
	}
}

func monsterProficiencyResultsToMonsterProficiencies(input []*monsterProficiency) []*entities.MonsterProficiency {
	out := make([]*entities.MonsterProficiency, len(input))
	for i, p := range input {
		out[i] = monsterProficiencyResultToMonsterProficiency(p)
	}

	return out
}

func monsterSensesResultToMonsterSenses(input *monsterSenses) *entities.MonsterSenses {
	if input == nil {
		return nil
	}

	return &entities.MonsterSenses{
		Blindsight:        input.Blindsight,
		Darkvision:        input.Darkvision,
		Tremorsense:       input.Tremorsense,
		Truesight:         input.Truesight,
		PassivePerception: input.PassivePerception,
	}
}

func monsterActionResultToMonsterAction(input *monsterAction) *entities.MonsterAction {
	if input == nil {
		return nil
	}

	return &entities.MonsterAction{
		Name:        input.Name,
		Description: input.Description,
		AttackBonus: input.AttackBonus,
		Damage:      damageResultsToDamage(input.Damage),
	}
}

func monsterActionResultsToMonsterActions(input []*monsterAction) []*entities.MonsterAction {
	out := make([]*entities.MonsterAction, len(input))
	for i, a := range input {
		out[i] = monsterActionResultToMonsterAction(a)
	}

	return out
}

func damageResultsToDamage(input []*damage) []*entities.Damage {
	out := make([]*entities.Damage, len(input))
	for i, d := range input {
		out[i] = damageResultToDamage(d)
	}

	return out
}

func spellCastingResultToSpellCasting(input *spellCasting) *entities.SpellCasting {
	if input == nil {
		return nil
	}

	return &entities.SpellCasting{
		CantripsKnown:    input.CantripsKnown,
		SpellsKnown:      input.SpellsKnown,
		SpellSlotsLevel1: input.SpellSlotsLevel1,
		SpellSlotsLevel2: input.SpellSlotsLevel2,
		SpellSlotsLevel3: input.SpellSlotsLevel3,
		SpellSlotsLevel4: input.SpellSlotsLevel4,
		SpellSlotsLevel5: input.SpellSlotsLevel5,
		SpellSlotsLevel6: input.SpellSlotsLevel6,
		SpellSlotsLevel7: input.SpellSlotsLevel7,
		SpellSlotsLevel8: input.SpellSlotsLevel8,
		SpellSlotsLevel9: input.SpellSlotsLevel9,
	}
}

func classSpecificResultToRangerSpecific(input *classSpecificResult) *entities.RangerSpecific {
	if input == nil {
		return nil
	}

	return &entities.RangerSpecific{
		FavoredEnemies: input.FavoredEnemies,
		FavoredTerrain: input.FavoredTerrain,
	}
}

func levelResultToClassSpecific(input *levelResult) entities.ClassSpecific {
	if input == nil {
		return nil
	}

	if input.Class == nil {
		return nil
	}

	switch input.Class.Index {
	case "ranger":
		return classSpecificResultToRangerSpecific(input.ClassSpecific)

	case "barbarian":
		return classSpecificResultToBarbarianSpecific(input.ClassSpecific)

	case "bard":
		return classSpecificResultToBardSpecific(input.ClassSpecific)

	case "cleric":
		return classSpecificResultToClericSpecific(input.ClassSpecific)

	case "druid":
		return classSpecificResultToDruidSpecific(input.ClassSpecific)

	case "fighter":
		return classSpecificResultToFighterSpecific(input.ClassSpecific)

	case "monk":
		return classSpecificResultToMonkSpecific(input.ClassSpecific)

	case "paladin":
		return classSpecificResultToPaladinSpecific(input.ClassSpecific)

	case "rogue":
		return classSpecificResultToRogueSpecific(input.ClassSpecific)

	case "sorcerer":
		return classSpecificResultToSorcererSpecific(input.ClassSpecific)

	case "warlock":
		return classSpecificResultToWarlockSpecific(input.ClassSpecific)

	case "wizard":
		return classSpecificResultToWizardSpecific(input.ClassSpecific)
	}

	return nil
}

func classSpecificResultToBarbarianSpecific(input *classSpecificResult) *entities.BarbarianSpecific {
	if input == nil {
		return nil
	}

	return &entities.BarbarianSpecific{
		RageCount:          input.RageCount,
		RageDamageBonus:    input.RageDamageBonus,
		BrutalCriticalDice: input.BrutalCriticalDice,
	}
}

func classSpecificResultToBardSpecific(input *classSpecificResult) *entities.BardSpecific {
	if input == nil {
		return nil
	}

	return &entities.BardSpecific{
		BardicInspirationDie: input.BardicInspirationDie,
		SongOfRestDie:        input.SongOfRestDie,
		MagicalSecretsMax5:   input.MagicalSecretsMax5,
		MagicalSecretsMax7:   input.MagicalSecretsMax7,
		MagicalSecretsMax9:   input.MagicalSecretsMax9,
	}
}

func classSpecificResultToClericSpecific(input *classSpecificResult) *entities.ClericSpecific {
	if input == nil {
		return nil
	}

	return &entities.ClericSpecific{
		ChannelDivinityCharges: input.ChannelDivinityCharges,
		DestroyUndeadCR:        input.DestroyUndeadCR,
	}
}

func classSpecificResultToDruidSpecific(input *classSpecificResult) *entities.DruidSpecific {
	if input == nil {
		return nil
	}

	return &entities.DruidSpecific{
		WildShapeMaxCR: input.WildShapeMaxCR,
		WildShapeSwim:  input.WildShapeSwim,
		WildShapeFly:   input.WildShapeFly,
	}
}

func classSpecificResultToFighterSpecific(input *classSpecificResult) *entities.FighterSpecific {
	if input == nil {
		return nil
	}

	return &entities.FighterSpecific{
		ActionSurges:    input.ActionSurges,
		IndomitableUses: input.IndomitableUses,
		ExtraAttacks:    input.ExtraAttacks,
	}
}

func classSpecificResultToMonkSpecific(input *classSpecificResult) *entities.MonkSpecific {
	if input == nil {
		return nil
	}

	return &entities.MonkSpecific{
		MartialArts:       martialArtsResultToMartialArts(input.MartialArts),
		KiPoints:          input.KiPoints,
		UnarmoredMovement: input.UnarmoredMovement,
	}
}

func martialArtsResultToMartialArts(input *martialArts) *entities.MartialArts {
	if input == nil {
		return nil
	}

	return &entities.MartialArts{
		DiceCount: input.DiceCount,
		DiceValue: input.DiceValue,
	}
}

func classSpecificResultToPaladinSpecific(input *classSpecificResult) *entities.PaladinSpecific {
	if input == nil {
		return nil
	}

	return &entities.PaladinSpecific{
		AuraRange: input.AuraRange,
	}
}

func sneakAttackResultToSneakAttack(input *sneakAttack) *entities.SneakAttack {
	if input == nil {
		return nil
	}

	return &entities.SneakAttack{
		DiceCount: input.DiceCount,
		DiceValue: input.DiceValue,
	}
}

func classSpecificResultToRogueSpecific(input *classSpecificResult) *entities.RogueSpecific {
	if input == nil {
		return nil
	}

	return &entities.RogueSpecific{
		SneakAttack: sneakAttackResultToSneakAttack(input.SneakAttack),
	}
}

func classSpecificResultToSorcererSpecific(input *classSpecificResult) *entities.SorcererSpecific {
	if input == nil {
		return nil
	}

	return &entities.SorcererSpecific{
		SorceryPoints:      input.SorceryPoints,
		MetamagicKnown:     input.MetamagicKnown,
		CreatingSpellSlots: creatingSpellSlotResultsToSpellSlots(input.CreatingSpellSlots),
	}
}

func creatingSpellSlotResultToSpellSlot(input *creatingSpellSlots) *entities.CreatingSpellSlots {
	if input == nil {
		return nil
	}

	return &entities.CreatingSpellSlots{
		SpellSlotLevel:   input.SpellSlotLevel,
		SorceryPointCost: input.SorceryPointCost,
	}
}

func creatingSpellSlotResultsToSpellSlots(input []*creatingSpellSlots) []*entities.CreatingSpellSlots {
	if input == nil {
		return nil
	}

	slots := make([]*entities.CreatingSpellSlots, len(input))
	for i := range input {
		slots[i] = creatingSpellSlotResultToSpellSlot(input[i])
	}

	return slots
}

func classSpecificResultToWarlockSpecific(input *classSpecificResult) *entities.WarlockSpecific {
	if input == nil {
		return nil
	}

	return &entities.WarlockSpecific{
		InvocationsKnown:    input.InvocationsKnown,
		MysticArcanumLevel6: input.MysticArcanumLevel6,
		MysticArcanumLevel7: input.MysticArcanumLevel7,
		MysticArcanumLevel8: input.MysticArcanumLevel8,
		MysticArcanumLevel9: input.MysticArcanumLevel9,
	}
}

func classSpecificResultToWizardSpecific(input *classSpecificResult) *entities.WizardSpecific {
	if input == nil {
		return nil
	}

	return &entities.WizardSpecific{
		ArcaneRecoveryLevels: input.ArcaneRecoveryLevels,
	}
}

func typeStringToProficiencyType(input string) entities.ProficiencyType {
	switch input {
	case "Armor":
		return entities.ProficiencyTypeArmor
	case "Weapons":
		return entities.ProficiencyTypeWeapon
	case "Artisan's Tools":
		return entities.ProficiencyTypeTool
	case "Musical Instruments":
		return entities.ProficiencyTypeInstrument
	case "Saving Throws":
		return entities.ProficiencyTypeSavingThrow
	case "Skills":
		return entities.ProficiencyTypeSkill
	default:
		return entities.ProficiencyTypeUnknown
	}
}

func referenceItemsToReferenceItems(input []*referenceItem) []*entities.ReferenceItem {
	if input == nil {
		return nil
	}

	items := make([]*entities.ReferenceItem, len(input))
	for i := range input {
		items[i] = referenceItemToReferenceItem(input[i])
	}

	return items
}

func referenceItemToReferenceItem(input *referenceItem) *entities.ReferenceItem {
	if input == nil {
		return nil
	}

	return &entities.ReferenceItem{
		Key:  input.Index,
		Name: input.Name,
		Type: urlToType(input.URL),
	}
}
