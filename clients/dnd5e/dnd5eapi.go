package dnd5e

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/fadedpez/dnd5e-api/entities"
)

const defaultBaserulzURL = "https://www.dnd5eapi.co/api/"

type dnd5eAPI struct {
	baserulzURL string
	client      httpIface
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

	if cfg.BaserulzURL == "" {
		cfg.BaserulzURL = defaultBaserulzURL
	}

	return &dnd5eAPI{client: cfg.Client, baserulzURL: cfg.BaserulzURL}, nil
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
		out[i] = referenceItemToReferenceItem(r)
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
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) listEquipmentByCategory(category string) ([]*referenceItem, error) {
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

	return response.Equipment, nil
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
		out[i] = referenceItemToReferenceItem(r)
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
		Proficiencies:            referenceItemsToReferenceItems(response.Proficiencies),
		SavingThrows:             referenceItemsToReferenceItems(response.SavingThrows),
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
				input.From.Options[idx] = opti choice" {
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
		url = baserulzURL + "spells"
	} else {
		url = baserulzURL + "spells?level=" + strconv.Itoa(*level)
	}

	resp, err := c.client.Get(url)
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

	return response.Results, nil
}

func (c *dnd5eAPI) doGetSpellsByClass(class string) ([]*referenceItem, error) {
	if class == "" {
		return nil, errors.New("class is empty")
	}

	url := baserulzURL + "classes/" + class + "/spells"

	resp, err := c.client.Get(url)
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

	return response.Results, nil
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
		SpellSchool:   referenceItemToReferenceItem(response.SpellSchool),
		SpellClasses:  referenceItemsToReferenceItems(response.SpellClasses),
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
		out[i] = referenceItemToReferenceItem(r)
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
	resp, err := c.client.Get(baserulzURL + "skills")
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
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetSkill(key string) (*entities.Skill, error) {
	resp, err := c.client.Get(baserulzURL + "skills/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
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
		Descricption: response.Description,
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
	url := baserulzURL + "monsters"
	
	// Add query parameters if provided
	if input != nil && input.ChallengeRating != nil {
		url = fmt.Sprintf("%s?challenge_rating=%g", url, *input.ChallengeRating)
	}
	
	resp, err := c.client.Get(url)
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
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetMonster(key string) (*entities.Monster, error) {
	resp, err := c.client.Get(baserulzURL + "monsters/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
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

	resp, err := c.client.Get(baserulzURL + "classes/" + key + "/levels/" + strconv.Itoa(level))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
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
	resp, err := c.client.Get(baserulzURL + "proficiencies/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
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
	resp, err := c.client.Get(baserulzURL + "damage-types")
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
		out[i] = referenceItemToReferenceItem(r)
	}

	return out, nil
}

func (c *dnd5eAPI) GetDamageType(key string) (*entities.DamageType, error) {
	resp, err := c.client.Get(baserulzURL + "damage-types/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
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
	resp, err := c.client.Get(baserulzURL + "equipment-categories/" + key)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()

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
