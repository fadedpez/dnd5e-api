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
		StartingProficiencyOptions: response.getStartingProficiencyChoice(),
		LanguageOptions:            response.getLanguageChoice(),
	}

	raceChoice := response.getStartingProficiencyChoice()
	if raceChoice == nil {
		return race, nil
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

func (c *dnd5eAPI) startingEquipmentCategoriesToOptionList(input []map[string]interface{}) ([]map[string]interface{}, error) {
	out := make([]map[string]interface{}, len(input))
	for i, v := range input {
		item, err := c.startingEquipmentCategoryToOptionList(v)
		if err != nil {
			return nil, err
		}
		out[i] = item
	}

	return out, nil
}

func (c *dnd5eAPI) startingEquipmentCategoryToOptionList(input map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for k, v := range input {
		if k == "from" {
			if item, ok := v.(map[string]interface{}); ok {
				if item["option_set_type"] == "equipment_category" {
					if category, ok := item["equipment_category"].(map[string]interface{}); ok {
						if name, ok := category["index"].(string); ok {
							equipment, err := c.listEquipmentByCategory(name)
							if err != nil {
								return nil, err
							}
							options := make([]interface{}, len(equipment))
							for idx, e := range equipment {
								options[idx] = map[string]interface{}{
									"option_type": "reference",
									"item": map[string]interface{}{
										"index": e.Key,
										"name":  e.Name,
										"url":   e.URL,
									},
								}

							}
							out[k] = map[string]interface{}{
								"option_set_type": "options_array",
								"options":         options,
							}
						}
					}
				} else if item["option_set_type"] == "options_array" {
					if options, ok := item["options"].([]interface{}); ok {
						for idx, optionItem := range options {
							if option, ok := optionItem.(map[string]interface{}); ok {
								if option["option_type"] == "choice" {
									if choice, ok := option["choice"].(map[string]interface{}); ok {
										newChoice, err := c.startingEquipmentCategoryToOptionList(choice)
										if err != nil {
											return nil, err
										}

										option["choice"] = newChoice
									}
								}
								options[idx] = option
							}

						}
						out[k] = item
					}

				}
			}

		} else {
			out[k] = v
		}
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

	startingEquipmentOption, err := c.startingEquipmentCategoriesToOptionList(response.StartingEquipmentOptions)
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
		ProficiencyChoices:       mapsToChoices(response.ProficiencyChoices),
		StartingEquipmentOptions: mapsToChoices(startingEquipmentOption),
	}

	return class, nil
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
