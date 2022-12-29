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
