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

	class := &entities.Class{
		Key:               response.Index,
		Name:              response.Name,
		HitDie:            response.HitDie,
		Proficiencies:     proficiencyResultsToProficiencies(response.Proficiencies),
		SavingThrows:      savingThrowResultsToSavingThrows(response.SavingThrows),
		StartingEquipment: startingEquipmentResultsToStartingEquipment(response.StartingEquipment),
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