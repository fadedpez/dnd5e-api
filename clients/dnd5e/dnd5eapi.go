package dnd5e

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/fadedpez/dnd5e-api/entities"
)

const baserulzURL = "https://www.dnd5eapi.co/api/"

type dnd5eAPI struct {
	client httpIface
}

type DND5eAPIConfig struct {
	Client httpIface
}

func NewDND5eAPI(cfg *DND5eAPIConfig) (Interface, error) {
	if cfg == nil {
		return nil, errors.New("cfg is required")
	}

	if cfg.Client == nil {
		return nil, errors.New("cfg.Client is required")
	}

	return &dnd5eAPI{client: cfg.Client}, nil
}

func (c *dnd5eAPI) ListRaces() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(baserulzURL + "races")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := listResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	out := make([]*entities.ReferenceItem, len(response.Results))
	for i, r := range response.Results {
		out[i] = listResultToRace(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetRace(key string) (*entities.Race, error) {
	resp, err := c.client.Get(baserulzURL + "races/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := raceResult{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	race := &entities.Race{
		Key:                        response.Index,
		Name:                       response.Name,
		Speed:                      response.Speed,
		AbilityBonuses:             abilityBonusResultsToAbilityBonuses(response.AbilityBonus),
		Languages:                  languageResultsToLanguages(response.Language),
		Traits:                     traitResultsToTraits(response.Trait),
		SubRaces:                   subRaceResultsToSubRaces(response.SubRaces),
		StartingProficiencies:      proficiencyResultsToProficiencies(response.StartingProficiencies),
		StartingProficiencyOptions: choiceResultToChoice(response.StartingProficiencyOptions),
		LanguageOptions:            choiceResultToChoice(response.LanguageOptions),
	}

	return race, nil
}

func (c *dnd5eAPI) ListEquipment() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(baserulzURL + "equipment")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := listResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	out := make([]*entities.ReferenceItem, len(response.Results))
	for i, r := range response.Results {
		out[i] = listResultToEquipment(r)
	}

	return out, nil
}

func (c *dnd5eAPI) listEquipmentByCategory(category string) ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(baserulzURL + "equipment-categories/" + category)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := equipmentListResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	out := make([]*entities.ReferenceItem, len(response.Equipment))
	for i, r := range response.Equipment {
		out[i] = listResultToEquipment(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetEquipment(key string) (EquipmentInterface, error) {
	resp, err := c.client.Get(baserulzURL + "equipment/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := equipmentResult{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	switch response.getCategoryKey() {
	case "weapon":
		weaponResponse := &weaponResult{}

		err = json.Unmarshal(responseBody, &weaponResponse)
		if err != nil {
			return nil, err
		}

		return weaponResultToWeapon(weaponResponse), nil

	case "armor":
		armorResponse := &armorResult{}

		err = json.Unmarshal(responseBody, &armorResponse)
		if err != nil {
			return nil, err
		}

		return armorResultToArmor(armorResponse), nil

	default:
		return equipmentResultToEquipment(&response), nil
	}
}

func (c *dnd5eAPI) ListClasses() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(baserulzURL + "classes")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := listResponse{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	out := make([]*entities.ReferenceItem, len(response.Results))
	for i, r := range response.Results {
		out[i] = listClassResultToClass(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetClass(key string) (*entities.Class, error) {
	resp, err := c.client.Get(baserulzURL + "classes/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := classResult{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	startingEquipmentOption, err := c.replaceEquipmentCategoryOptionSetTypesToOptionsArrays(response.StartingEquipmentOptions)
	if err != nil {
		return nil, err
	}

	class := &entities.Class{
		Key:                      response.Index,
		Name:                     response.Name,
		HitDie:                   response.HitDie,
		Proficiencies:            proficiencyResultsToProficiencies(response.Proficiencies),
		SavingThrows:             savingThrowResultsToSavingThrows(response.SavingThrows),
		StartingEquipment:        startingEquipmentResultsToStartingEquipment(response.StartingEquipment),
		ProficiencyChoices:       choiceResultsToChoices(response.ProficiencyChoices),
		StartingEquipmentOptions: startingEquipmentOption,
	}

	return class, nil
}

func (c *dnd5eAPI) replaceEquipmentCategoryOptionSetTypesToOptionsArrays(input []*choiceResult) ([]*entities.ChoiceOption, error) {
	out := make([]*entities.ChoiceOption, len(input))
	for i, item := range input { // item is a choice
		newChoice, err := c.replaceEquipmentCategoryOptionSetTypeToOptionsArray(item)
		if err != nil {
			return nil, err
		}
		out[i] = newChoice.toEntity()
	}
	return out, nil
}

func (c *dnd5eAPI) replaceEquipmentCategoryOptionSetTypeToOptionsArray(input *choiceResult) (*choiceResult, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if input.Type != "equipment" {
		return input, nil
	}

	if input.From == nil {
		return input, nil
	}

	if input.From.OptionSetType == "options_array" {
		for idx, option := range input.From.Options {
			if option.OptionType == "choice" {
				newChoice, err := c.replaceEquipmentCategoryOptionSetTypeToOptionsArray(option.Choice)
				if err != nil {
					return nil, err
				}
				option.Choice = newChoice
				input.From.Options[idx] = option
			} else if option.OptionType == "multiple" {
				for idx2, multiple := range option.Items {
					if multiple.OptionType == "choice" {
						newChoice, err := c.replaceEquipmentCategoryOptionSetTypeToOptionsArray(multiple.Choice)
						if err != nil {
							return nil, err
						}
						multiple.Choice = newChoice
						option.Items[idx2] = multiple //TODO: refactor/rename idx2?
					}

				}
				input.From.Options[idx] = option
			}
		}

		return input, nil
	}

	if input.From.OptionSetType != "equipment_category" {
		return input, nil
	}

	//TODO: should we return an error?
	if input.From.EquipmentCategory == nil {
		return input, nil
	}

	equipment, err := c.listEquipmentByCategory(input.From.EquipmentCategory.Index)
	if err != nil {
		return nil, err
	}

	input.From.OptionSetType = "options_array"

	options := make([]*option, len(equipment))
	for i, e := range equipment {
		options[i] = &option{
			OptionType: "reference",
			Item: &referenceItem{
				Name:  e.Name,
				URL:   fmt.Sprintf("/api/equipment/%s", e.Key),
				Index: e.Key,
			},
		}
	}

	input.From.Options = options

	return input, nil
}

func (c *dnd5eAPI) ListSpells() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(baserulzURL + "spells")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := listResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	out := make([]*entities.ReferenceItem, len(response.Results))
	for i, r := range response.Results {
		out[i] = listResultToSpell(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetSpell(key string) (*entities.Spell, error) {
	resp, err := c.client.Get(baserulzURL + "spells/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := spellResult{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	spell := &entities.Spell{
		Key:           response.Index,
		Name:          response.Name,
		Range:         response.Range,
		Ritual:        response.Ritual,
		Duration:      response.Duration,
		Concentration: response.Concentration,
		CastingTime:   response.CastingTime,
		SpellLevel:    response.SpellLevel,
		SpellDamage:   spellDamageResultToSpellDamage(response.SpellDamage),
		DC:            dcResultToDC(response.DC),
		AreaOfEffect:  areaOfEffectResultToAreaOfEffect(response.AreaOfEffect),
		SpellSchool:   spellSchoolResultToSpellSchool(response.SpellSchool),
		SpellClasses:  spellClassResultsToSpellClasses(response.SpellClasses),
	}

	return spell, nil
}

func (c *dnd5eAPI) ListFeatures() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(baserulzURL + "features")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := listResponse{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, err
	}

	out := make([]*entities.ReferenceItem, len(response.Results))
	for i, r := range response.Results {
		out[i] = listResultToFeature(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetFeature(key string) (*entities.Feature, error) {
	resp, err := c.client.Get(baserulzURL + "features/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
	response := featureResult{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return nil, err
	}

	feature := &entities.Feature{
		Key:             response.Index,
		Name:            response.Name,
		Level:           response.Level,
		Class:           featureClassResultToClass(response.Class),
		FeatureSpecific: choiceResultToChoice(response.SubFeatureOptions),
	}

	return feature, nil
}
