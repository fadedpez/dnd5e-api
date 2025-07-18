package dnd5e

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/fadedpez/dnd5e-api/entities"
)

const (
	baserulzURL    = "https://www.dnd5eapi.co/api/"
	httpStatusOK   = 200
)

type dnd5eAPI struct {
	client  httpIface
	baseURL string
	mu      sync.RWMutex
}

type DND5eAPIConfig struct {
	Client  httpIface
	BaseURL string
}

func NewDND5eAPI(cfg *DND5eAPIConfig) (Interface, error) {
	if cfg == nil {
		return nil, errors.New("cfg is required")
	}

	if cfg.Client == nil {
		return nil, errors.New("cfg.Client is required")
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = baserulzURL
	}

	return &dnd5eAPI{
		client:  cfg.Client,
		baseURL: baseURL,
	}, nil
}

func (c *dnd5eAPI) getBaseURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.baseURL
}

// newHTTPStatusError creates a standardized error for unexpected HTTP status codes
func newHTTPStatusError(statusCode int) error {
	return fmt.Errorf("unexpected status code: %d", statusCode)
}

func (c *dnd5eAPI) ListRaces() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(c.getBaseURL() + "races")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := listResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	out := make([]*entities.ReferenceItem, len(response.Results))
	for i, r := range response.Results {
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetRace(key string) (*entities.Race, error) {
	resp, err := c.client.Get(c.getBaseURL() + "races/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := raceResult{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	race := &entities.Race{
		Key:                        response.Index,
		Name:                       response.Name,
		Speed:                      response.Speed,
		Size:                       response.Size,
		SizeDescription:            response.SizeDescription,
		AbilityBonuses:             abilityBonusResultsToAbilityBonuses(response.AbilityBonus),
		Languages:                  referenceItemsToReferenceItems(response.Language),
		Traits:                     referenceItemsToReferenceItems(response.Trait),
		SubRaces:                   referenceItemsToReferenceItems(response.SubRaces),
		StartingProficiencies:      referenceItemsToReferenceItems(response.StartingProficiencies),
		StartingProficiencyOptions: choiceResultToChoice(response.StartingProficiencyOptions),
		LanguageOptions:            choiceResultToChoice(response.LanguageOptions),
	}

	return race, nil
}

func (c *dnd5eAPI) ListEquipment() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(c.getBaseURL() + "equipment")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := listResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	out := make([]*entities.ReferenceItem, len(response.Results))
	for i, r := range response.Results {
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) listEquipmentByCategory(category string) ([]*referenceItem, error) {
	resp, err := c.client.Get(c.getBaseURL() + "equipment-categories/" + category)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := equipmentListResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Equipment, nil
}

func (c *dnd5eAPI) GetEquipment(key string) (EquipmentInterface, error) {
	resp, err := c.client.Get(c.getBaseURL() + "equipment/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
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
	resp, err := c.client.Get(c.getBaseURL() + "classes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
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
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetClass(key string) (*entities.Class, error) {
	resp, err := c.client.Get(c.getBaseURL() + "classes/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
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

	armorProfs, weaponProfs, toolProfs := categorizeProficiencies(response.Proficiencies)
	
	class := &entities.Class{
		Key:                      response.Index,
		Name:                     response.Name,
		HitDie:                   response.HitDie,
		Proficiencies:            referenceItemsToReferenceItems(response.Proficiencies),
		SavingThrows:             referenceItemsToReferenceItems(response.SavingThrows),
		StartingEquipment:        startingEquipmentResultsToStartingEquipment(response.StartingEquipment),
		ProficiencyChoices:       choiceResultsToChoices(response.ProficiencyChoices),
		StartingEquipmentOptions: startingEquipmentOption,
		PrimaryAbilities:         extractPrimaryAbilities(response.MultiClassing),
		Description:              getClassDescription(response.Index),
		ArmorProficiencies:       armorProfs,
		WeaponProficiencies:      weaponProfs,
		ToolProficiencies:        toolProfs,
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
			Item:       e,
		}
	}

	input.From.Options = options

	return input, nil
}

type ListSpellsInput struct {
	Level *int
	Class string
}

func (c *dnd5eAPI) ListSpells(input *ListSpellsInput) ([]*entities.ReferenceItem, error) {
	if input == nil {
		return nil, errors.New("input is nil")
	}

	if input.Class == "" {
		levelList, err := c.doGetSpellsByLevel(input.Level)
		if err != nil {
			return nil, err
		}

		levelOut := make([]*entities.ReferenceItem, len(levelList))
		for i, r := range levelList {
			levelOut[i] = referenceItemToReferenceItem(r)
		}

		return levelOut, nil
	}

	if input.Level == nil {
		classList, err := c.doGetSpellsByClass(input.Class)
		if err != nil {
			return nil, err
		}

		classOut := make([]*entities.ReferenceItem, len(classList))
		for i, r := range classList {
			classOut[i] = referenceItemToReferenceItem(r)
		}

		return classOut, nil
	}

	levelList, err := c.doGetSpellsByLevel(input.Level)
	if err != nil {
		return nil, err
	}

	levelMap := make(map[string]bool)
	for _, r := range levelList {
		levelMap[r.Index] = true
	}

	classList, err := c.doGetSpellsByClass(input.Class)
	if err != nil {
		return nil, err
	}

	out := make([]*entities.ReferenceItem, 0)
	for _, r := range classList {
		if levelMap[r.Index] {
			out = append(out, referenceItemToReferenceItem(r))
		}
	}

	return out, nil
}

func (c *dnd5eAPI) doGetSpellsByLevel(level *int) ([]*referenceItem, error) {
	var url string
	if level == nil {
		url = c.getBaseURL() + "spells"
	} else {
		url = c.getBaseURL() + "spells?level=" + strconv.Itoa(*level)
	}

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := listResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Results, nil
}

func (c *dnd5eAPI) doGetSpellsByClass(class string) ([]*referenceItem, error) {
	if class == "" {
		return nil, errors.New("class is empty")
	}

	url := c.getBaseURL() + "classes/" + class + "/spells"

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := listResponse{}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response.Results, nil
}

func (c *dnd5eAPI) GetSpell(key string) (*entities.Spell, error) {
	resp, err := c.client.Get(c.getBaseURL() + "spells/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
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
		SpellSchool:   referenceItemToReferenceItem(response.SpellSchool),
		SpellClasses:  referenceItemsToReferenceItems(response.SpellClasses),
	}

	return spell, nil
}

func (c *dnd5eAPI) ListFeatures() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(c.getBaseURL() + "features")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
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
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetFeature(key string) (*entities.Feature, error) {
	resp, err := c.client.Get(c.getBaseURL() + "features/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
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
		Key:   response.Index,
		Name:  response.Name,
		Level: response.Level, //TODO: add prerequisites?
		Class: referenceItemToReferenceItem(response.Class),
	}

	if response.FeatureSpecific != nil {
		feature.FeatureSpecific = &entities.SubFeatureOption{
			SubFeatureOptions: choiceResultToChoice(response.FeatureSpecific.SubFeatureOptions),
		}
	}

	return feature, nil
}

func (c *dnd5eAPI) ListSkills() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(c.getBaseURL() + "skills")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
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
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetSkill(key string) (*entities.Skill, error) {
	resp, err := c.client.Get(c.getBaseURL() + "skills/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := skillResult{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return nil, err
	}

	skill := &entities.Skill{
		Key:          response.Index,
		Name:         response.Name,
		Description: response.Description,
		AbilityScore: referenceItemToReferenceItem(response.AbilityScore),
		Type:         urlToType(response.URL),
	}

	return skill, nil
}

type ListMonstersInput struct {
	ChallengeRating *float64
}

func (c *dnd5eAPI) ListMonsters() ([]*entities.ReferenceItem, error) {
	return c.ListMonstersWithFilter(nil)
}

func (c *dnd5eAPI) ListMonstersWithFilter(input *ListMonstersInput) ([]*entities.ReferenceItem, error) {
	url := c.getBaseURL() + "monsters"
	
	// Add query parameters if provided
	if input != nil && input.ChallengeRating != nil {
		url = fmt.Sprintf("%s?challenge_rating=%g", url, *input.ChallengeRating)
	}
	
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
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
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetMonster(key string) (*entities.Monster, error) {
	resp, err := c.client.Get(c.getBaseURL() + "monsters/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := monsterResult{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return nil, err
	}

	monster := &entities.Monster{
		Key:                   response.Index,
		Name:                  response.Name,
		Size:                  response.Size,
		Type:                  response.Type,
		Alignment:             response.Alignment,
		ArmorClass:            monsterArmorClassToValue(response.ArmorClass),
		HitPoints:             response.HitPoints,
		HitDice:               response.HitDice,
		Speed:                 monsterSpeedResultToSpeed(response.Speed),
		Strength:              response.Strength,
		Dexterity:             response.Dexterity,
		Constitution:          response.Constitution,
		Intelligence:          response.Intelligence,
		Wisdom:                response.Wisdom,
		Charisma:              response.Charisma,
		Proficiencies:         monsterProficiencyResultsToMonsterProficiencies(response.Proficiencies),
		DamageVulnerabilities: response.DamageVulnerabilities,
		DamageResistances:     response.DamageResistances,
		DamageImmunities:      response.DamageImmunities,
		ConditionImmunities:   referenceItemsToReferenceItems(response.ConditionImmunities),
		MonsterSenses:         monsterSensesResultToMonsterSenses(response.Senses),
		Languages:             response.Languages,
		ChallengeRating:       response.ChallengeRating,
		XP:                    response.XP,
		MonsterActions:        monsterActionResultsToMonsterActions(response.MonsterActions),
		MonsterImageURL:       response.MonsterImageURL,
	}

	return monster, nil
}

func monsterArmorClassToValue(input []*monsterArmorClass) int {
	if input == nil {
		return 0
	}

	if len(input) == 0 {
		return 0
	}

	sum := 0
	for _, ac := range input {
		sum += ac.Value
	}

	return sum
}
func (c *dnd5eAPI) GetClassLevel(key string, level int) (*entities.Level, error) {
	if key == "" {
		return nil, errors.New("key is required")
	}

	if level == 0 {
		return nil, errors.New("level is required")
	}

	resp, err := c.client.Get(c.getBaseURL() + "classes/" + key + "/levels/" + strconv.Itoa(level))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := &levelResult{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, response)

	if err != nil {
		return nil, err
	}

	classLevel := &entities.Level{
		Level:               response.Level,
		AbilityScoreBonuses: response.AbilityScoreBonuses,
		ProfBonus:           response.ProfBonus,
		Features:            referenceItemsToReferenceItems(response.Features),
		SpellCasting:        spellCastingResultToSpellCasting(response.SpellCasting),
		ClassSpecific:       levelResultToClassSpecific(response),
		Key:                 response.Index,
		Class:               referenceItemToReferenceItem(response.Class),
	}

	return classLevel, nil
}

func (c *dnd5eAPI) GetProficiency(key string) (*entities.Proficiency, error) {
	resp, err := c.client.Get(c.getBaseURL() + "proficiencies/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)

	response := proficiencyResult{}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return nil, err
	}

	proficiency := &entities.Proficiency{
		Key:       response.Index,
		Name:      response.Name,
		Type:      typeStringToProficiencyType(response.Type),
		Reference: referenceItemToReferenceItem(response.Reference),
	}

	return proficiency, nil
}

func (c *dnd5eAPI) ListDamageTypes() ([]*entities.ReferenceItem, error) {
	resp, err := c.client.Get(c.getBaseURL() + "damage-types")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
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
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetDamageType(key string) (*entities.DamageType, error) {
	resp, err := c.client.Get(c.getBaseURL() + "damage-types/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}
	response := damageTypeResult{}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(responseBody, &response)

	if err != nil {
		return nil, err
	}

	damageType := &entities.DamageType{
		Key:         response.Index,
		Name:        response.Name,
		Type:        urlToType(response.URL),
		Description: response.Description,
	}

	return damageType, nil
}

func (c *dnd5eAPI) GetEquipmentCategory(key string) (*entities.EquipmentCategory, error) {
	resp, err := c.client.Get(c.getBaseURL() + "equipment-categories/" + key)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		return nil, newHTTPStatusError(resp.StatusCode)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	category := &entities.EquipmentCategory{}
	err = json.Unmarshal(responseBody, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (c *dnd5eAPI) ListBackgrounds() ([]*entities.ReferenceItem, error) {
	// Try to get from API first
	resp, err := c.client.Get(c.getBaseURL() + "backgrounds")
	if err != nil {
		// If API fails, return hardcoded backgrounds
		return getHardcodedBackgrounds(), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		// If API returns error, return hardcoded backgrounds
		return getHardcodedBackgrounds(), nil
	}
	
	response := listResponse{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		// If parsing fails, return hardcoded backgrounds
		return getHardcodedBackgrounds(), nil
	}

	// Merge API results with hardcoded backgrounds
	apiBackgrounds := make([]*entities.ReferenceItem, len(response.Results))
	for i, r := range response.Results {
		apiBackgrounds[i] = referenceItemToReferenceItem(r)
	}
	
	// Get hardcoded backgrounds and merge, avoiding duplicates
	hardcodedBackgrounds := getHardcodedBackgrounds()
	merged := make([]*entities.ReferenceItem, 0, len(apiBackgrounds)+len(hardcodedBackgrounds))
	merged = append(merged, apiBackgrounds...)
	
	// Add hardcoded backgrounds that aren't in API results
	apiKeys := make(map[string]bool)
	for _, bg := range apiBackgrounds {
		apiKeys[bg.Key] = true
	}
	
	for _, bg := range hardcodedBackgrounds {
		if !apiKeys[bg.Key] {
			merged = append(merged, bg)
		}
	}

	return merged, nil
}

func (c *dnd5eAPI) GetBackground(key string) (*entities.Background, error) {
	// Try to get from API first
	resp, err := c.client.Get(c.getBaseURL() + "backgrounds/" + key)
	if err != nil {
		// If API fails, try hardcoded background
		return getHardcodedBackground(key)
	}
	defer resp.Body.Close()

	if resp.StatusCode != httpStatusOK {
		// If API returns 404 or other error, try hardcoded background
		return getHardcodedBackground(key)
	}
	
	response := backgroundResult{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		// If parsing fails, try hardcoded background
		return getHardcodedBackground(key)
	}

	background := &entities.Background{
		Key:                        response.Index,
		Name:                       response.Name,
		SkillProficiencies:         referenceItemsToReferenceItems(response.StartingProficiencies),
		LanguageOptions:            choiceResultToChoice(response.LanguageOptions),
		StartingEquipment:          startingEquipmentResultsToStartingEquipment(response.StartingEquipment),
		StartingEquipmentOptions:   choiceResultsToChoices(response.StartingEquipmentOptions),
		Feature:                    backgroundFeatureResultToBackgroundFeature(response.Feature),
		PersonalityTraits:          choiceResultToChoice(response.PersonalityTraits),
		Ideals:                     choiceResultToChoice(response.Ideals),
		Bonds:                      choiceResultToChoice(response.Bonds),
		Flaws:                      choiceResultToChoice(response.Flaws),
	}

	return background, nil
}

func extractPrimaryAbilities(multiclassing *multiClassing) []*entities.ReferenceItem {
	if multiclassing == nil || multiclassing.Prerequisites == nil {
		return nil
	}

	primaryAbilities := make([]*entities.ReferenceItem, 0, len(multiclassing.Prerequisites))
	for _, prereq := range multiclassing.Prerequisites {
		if prereq.AbilityScore != nil {
			primaryAbilities = append(primaryAbilities, referenceItemToReferenceItem(prereq.AbilityScore))
		}
	}

	return primaryAbilities
}

func getClassDescription(key string) string {
	descriptions := map[string]string{
		"barbarian":  "A fierce warrior of primitive background who can enter a battle rage",
		"bard":       "A master of song, speech, and the magic they contain",
		"cleric":     "A priestly champion who wields divine magic in service of a higher power",
		"druid":      "A priest of nature, wielding elemental forces and transformative magic",
		"fighter":    "A master of martial combat, skilled with a variety of weapons and armor",
		"monk":       "A master of martial arts, harnessing inner power through discipline",
		"paladin":    "A holy warrior bound to a sacred oath, wielding divine magic",
		"ranger":     "A warrior of the wilderness, skilled in tracking, survival, and combat",
		"rogue":      "A scoundrel who uses stealth and trickery to achieve their goals",
		"sorcerer":   "A spellcaster who draws on inherent magic from a gift or bloodline",
		"warlock":    "A wielder of magic derived from a bargain with an extraplanar entity",
		"wizard":     "A scholarly magic-user capable of manipulating structures of reality",
	}
	
	if description, exists := descriptions[key]; exists {
		return description
	}
	
	return ""
}

// getHardcodedBackgrounds returns a list of standard D&D 5e backgrounds
func getHardcodedBackgrounds() []*entities.ReferenceItem {
	return []*entities.ReferenceItem{
		{Key: "acolyte", Name: "Acolyte"},
		{Key: "criminal", Name: "Criminal"},
		{Key: "folk-hero", Name: "Folk Hero"},
		{Key: "noble", Name: "Noble"},
		{Key: "sage", Name: "Sage"},
		{Key: "soldier", Name: "Soldier"},
		{Key: "charlatan", Name: "Charlatan"},
		{Key: "entertainer", Name: "Entertainer"},
		{Key: "guild-artisan", Name: "Guild Artisan"},
		{Key: "hermit", Name: "Hermit"},
		{Key: "outlander", Name: "Outlander"},
		{Key: "sailor", Name: "Sailor"},
	}
}

// getHardcodedBackground returns detailed background data for standard D&D 5e backgrounds
func getHardcodedBackground(key string) (*entities.Background, error) {
	backgrounds := getHardcodedBackgroundData()
	if bg, exists := backgrounds[key]; exists {
		return bg, nil
	}
	return nil, fmt.Errorf("background not found: %s", key)
}

// getHardcodedBackgroundData returns detailed background data
func getHardcodedBackgroundData() map[string]*entities.Background {
	return map[string]*entities.Background{
		"criminal": {
			Key:  "criminal",
			Name: "Criminal",
			SkillProficiencies: []*entities.ReferenceItem{
				{Key: "skill-deception", Name: "Skill: Deception"},
				{Key: "skill-stealth", Name: "Skill: Stealth"},
			},
			Feature: &entities.BackgroundFeature{
				Name: "Criminal Contact",
				Description: "You have a reliable and trustworthy contact who acts as your liaison to a network of other criminals.",
			},
		},
		"folk-hero": {
			Key:  "folk-hero",
			Name: "Folk Hero",
			SkillProficiencies: []*entities.ReferenceItem{
				{Key: "skill-animal-handling", Name: "Skill: Animal Handling"},
				{Key: "skill-survival", Name: "Skill: Survival"},
			},
			Feature: &entities.BackgroundFeature{
				Name: "Rustic Hospitality",
				Description: "Since you come from the ranks of the common folk, you fit in among them with ease.",
			},
		},
		"sage": {
			Key:  "sage",
			Name: "Sage",
			SkillProficiencies: []*entities.ReferenceItem{
				{Key: "skill-arcana", Name: "Skill: Arcana"},
				{Key: "skill-history", Name: "Skill: History"},
			},
			Feature: &entities.BackgroundFeature{
				Name: "Researcher",
				Description: "When you attempt to learn or recall a piece of lore, if you do not know that information, you often know where and from whom you can obtain it.",
			},
		},
		"soldier": {
			Key:  "soldier",
			Name: "Soldier",
			SkillProficiencies: []*entities.ReferenceItem{
				{Key: "skill-athletics", Name: "Skill: Athletics"},
				{Key: "skill-intimidation", Name: "Skill: Intimidation"},
			},
			Feature: &entities.BackgroundFeature{
				Name: "Military Rank",
				Description: "You have a military rank from your career as a soldier. Soldiers loyal to your former military organization still recognize your authority and influence.",
			},
		},
		"noble": {
			Key:  "noble",
			Name: "Noble",
			SkillProficiencies: []*entities.ReferenceItem{
				{Key: "skill-history", Name: "Skill: History"},
				{Key: "skill-persuasion", Name: "Skill: Persuasion"},
			},
			Feature: &entities.BackgroundFeature{
				Name: "Position of Privilege",
				Description: "Thanks to your noble birth, people are inclined to think the best of you.",
			},
		},
		"charlatan": {
			Key:  "charlatan",
			Name: "Charlatan",
			SkillProficiencies: []*entities.ReferenceItem{
				{Key: "skill-deception", Name: "Skill: Deception"},
				{Key: "skill-sleight-of-hand", Name: "Skill: Sleight of Hand"},
			},
			Feature: &entities.BackgroundFeature{
				Name: "False Identity",
				Description: "You have created a second identity that includes documentation, established acquaintances, and disguises.",
			},
		},
	}
}

func categorizeProficiencies(proficiencies []*referenceItem) (armor, weapon, tool []*entities.ReferenceItem) {
	armorProficiencies := make([]*entities.ReferenceItem, 0)
	weaponProficiencies := make([]*entities.ReferenceItem, 0)
	toolProficiencies := make([]*entities.ReferenceItem, 0)

	for _, prof := range proficiencies {
		if prof == nil {
			continue
		}
		
		switch {
		case isArmorProficiency(prof.Index):
			armorProficiencies = append(armorProficiencies, referenceItemToReferenceItem(prof))
		case isWeaponProficiency(prof.Index):
			weaponProficiencies = append(weaponProficiencies, referenceItemToReferenceItem(prof))
		case isToolProficiency(prof.Index):
			toolProficiencies = append(toolProficiencies, referenceItemToReferenceItem(prof))
		}
	}

	return armorProficiencies, weaponProficiencies, toolProficiencies
}

func isArmorProficiency(index string) bool {
	armorProficiencies := map[string]bool{
		"light-armor":  true,
		"medium-armor": true,
		"heavy-armor":  true,
		"shields":      true,
		"all-armor":    true,
	}
	return armorProficiencies[index]
}

func isWeaponProficiency(index string) bool {
	weaponProficiencies := map[string]bool{
		"simple-weapons":  true,
		"martial-weapons": true,
	}
	return weaponProficiencies[index]
}

const savingThrowPrefix = "saving-throw"

func isToolProficiency(index string) bool {
	// Tools are proficiencies that are not armor, weapons, or saving throws
	// This handles various tool types like "smiths-tools", "thieves-tools", etc.
	// as well as any unknown proficiency types that may be added in the future
	return !isArmorProficiency(index) && 
		   !isWeaponProficiency(index) && 
		   !isSavingThrowProficiency(index)
}

func isSavingThrowProficiency(index string) bool {
	return strings.HasPrefix(index, savingThrowPrefix)
}
