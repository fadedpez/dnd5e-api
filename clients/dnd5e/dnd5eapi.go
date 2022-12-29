package dnd5e

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fadedpez/dnd5e-api/entities"
)

const baserulzURL = "https://www.dnd5eapi.co/api/"

type httpIface interface {
	Get(url string) (*http.Response, error)
}

type dnd5eAPI struct {
	client httpIface
}

type DND5eAPIConfig struct {
	Client httpIface
}

type listResult struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

type raceResult struct {
	Index                      string                 `json:"index"`
	Name                       string                 `json:"name"`
	Speed                      int                    `json:"speed"`
	AbilityBonus               []*abilityBonus        `json:"ability_bonuses"`
	Language                   []*language            `json:"languages"`
	Trait                      []*trait               `json:"traits"`
	SubRaces                   []*subRace             `json:"subraces"`
	StartingProficiencies      []*proficiency         `json:"starting_proficiencies"`
	StartingProficiencyOptions map[string]interface{} `json:"starting_proficiency_options"`
	LanguageOptions            map[string]interface{} `json:"language_options"`
}

func (r *raceResult) getStartingProficiencyChoice() *entities.Choice {
	if r.StartingProficiencyOptions == nil {
		return nil
	}

	out := &entities.Choice{
		Choose:    int(r.StartingProficiencyOptions["choose"].(float64)),
		Type:      r.StartingProficiencyOptions["type"].(string),
		OptionSet: mapToOptionSet(r.StartingProficiencyOptions["from"].(map[string]interface{})),
	}

	return out
}

func (r *raceResult) getLanguageChoice() *entities.Choice {
	if r.LanguageOptions == nil {
		return nil
	}

	out := &entities.Choice{
		Choose:    int(r.LanguageOptions["choose"].(float64)),
		Type:      r.LanguageOptions["type"].(string),
		OptionSet: mapToOptionSet(r.LanguageOptions["from"].(map[string]interface{})),
	}

	return out
}

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

type abilityBonus struct {
	AbilityScore *listResult `json:"ability_score"`
	Bonus        int         `json:"bonus"`
}

type language struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type trait struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type subRace struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type proficiency struct {
	Index string `json:"index"`
	Name  string `json:"name"`
}

type listResponse struct {
	Count   int           `json:"count"`
	Results []*listResult `json:"results"`
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

func (c *dnd5eAPI) ListRaces() ([]*entities.Race, error) {
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

	out := make([]*entities.Race, len(response.Results))
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
