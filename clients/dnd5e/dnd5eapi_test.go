package dnd5e

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/fadedpez/dnd5e-api/entities"

	"github.com/stretchr/testify/assert"
)

func TestNewDND5eAPI(t *testing.T) {
	t.Run("cfg is required", func(t *testing.T) {
		_, err := NewDND5eAPI(nil)
		if err == nil {
			t.Error("expected error, got nil")
		}

		assert.Equal(t, err.Error(), "cfg is required")
	})

	t.Run("cfg.Client is required", func(t *testing.T) {
		_, err := NewDND5eAPI(&DND5eAPIConfig{})
		if err == nil {
			t.Error("expected error, got nil")
		}

		assert.Equal(t, err.Error(), "cfg.Client is required")
	})

	t.Run("returns dnd5eAPI", func(t *testing.T) {
		actual, err := NewDND5eAPI(&DND5eAPIConfig{Client: &mockHTTPClient{}})

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.NotNil(t, actual.(*dnd5eAPI).client)
	})
}

func TestDND5eAPI_ListRaces(t *testing.T) {

	t.Run("returns error if http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"races").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListRaces()

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "http.Get failed")
	})

	t.Run("returns error if json.Decode fails", func(t *testing.T) {
		invalidJSONBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"count": 1, "results": [{"index": "human", "name": "Human", "url": "https://www.dnd5eapi"`)))
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"races").Return(&http.Response{
			StatusCode: 200,
			Body:       invalidJSONBody,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListRaces()

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected EOF", err.Error())
	})

	t.Run("returns error if response status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"races").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListRaces()

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("returns valid list of races", func(t *testing.T) {
		validJSONBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"count": 1, "results": [{"index": "human", "name": "Human", "url": "https://www.dnd5eapi"}]}`)))
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"races").Return(&http.Response{
			StatusCode: 200,
			Body:       validJSONBody,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		actual, err := dnd5eAPI.ListRaces()

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, 1, len(actual))
		assert.Equal(t, "human", actual[0].Key)
		assert.Equal(t, "Human", actual[0].Name)
	})
}

func TestDND5eAPI_GetRace(t *testing.T) {

	t.Run("returns error if http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"races/human").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetRace("human")

		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), "http.Get failed")
	})

	t.Run("returns error if json.Decode fails", func(t *testing.T) {
		invalidJSONBody := ioutil.NopCloser(bytes.NewReader([]byte(`{"index": "human", "name": "Human", "url": "https://www.dnd5eapi"`)))
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"races/human").Return(&http.Response{
			StatusCode: 200,
			Body:       invalidJSONBody,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetRace("human")

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected EOF", err.Error())
	})

	t.Run("returns error if response status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"races/human").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetRace("human")

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("returns a race", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/races/human.json")
		raceFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"races/human").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(raceFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		actual, err := dnd5eAPI.GetRace("human")

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, "human", actual.Key)
		assert.Equal(t, "Human", actual.Name)
		assert.Equal(t, 30, actual.Speed)
		assert.Equal(t, 6, len(actual.AbilityBonuses))
		assert.Equal(t, "str", actual.AbilityBonuses[0].AbilityScore.Key)
		assert.Equal(t, "STR", actual.AbilityBonuses[0].AbilityScore.Name)
		assert.Equal(t, 1, actual.AbilityBonuses[0].Bonus)
		assert.Equal(t, "dex", actual.AbilityBonuses[1].AbilityScore.Key)
		assert.Equal(t, "DEX", actual.AbilityBonuses[1].AbilityScore.Name)
		assert.Equal(t, 1, actual.AbilityBonuses[1].Bonus)
		assert.Equal(t, actual.Languages[0].Key, "common")
		assert.Equal(t, 1, actual.LanguageOptions.Choose)
		assert.Equal(t, "languages", actual.LanguageOptions.Type)
		assert.Equal(t, entities.OptionSetTypeArray, actual.LanguageOptions.OptionSet.GetType())
		assert.Equal(t, 15, len(actual.LanguageOptions.OptionSet.(*entities.OptionsArrayOptionSet).Options))
		assert.Equal(t, "dwarvish", actual.LanguageOptions.OptionSet.(*entities.OptionsArrayOptionSet).Options[0].(*entities.ReferenceOption).Reference.Key)
		assert.Equal(t, "elvish", actual.LanguageOptions.OptionSet.(*entities.OptionsArrayOptionSet).Options[1].(*entities.ReferenceOption).Reference.Key)
	})

	t.Run("returns a trait", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/races/elf.json")
		raceFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"races/elf").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(raceFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		actual, err := dnd5eAPI.GetRace("elf")

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, "darkvision", actual.Traits[0].Key)
		assert.Equal(t, "fey-ancestry", actual.Traits[1].Key)
	})

	t.Run("returns a subrace", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/races/elf.json")
		raceFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"races/elf").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(raceFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		actual, err := dnd5eAPI.GetRace("elf")

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, "high-elf", actual.SubRaces[0].Key)
	})

	t.Run("it returns starting proficiencies", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/races/dwarf.json")
		raceFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"races/dwarf").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(raceFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		actual, err := dnd5eAPI.GetRace("dwarf")

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, "battleaxes", actual.StartingProficiencies[0].Key)
		assert.Equal(t, "handaxes", actual.StartingProficiencies[1].Key)
		assert.Equal(t, 1, actual.StartingProficiencyOptions.Choose)
		assert.Equal(t, "proficiencies", actual.StartingProficiencyOptions.Type)
		assert.Equal(t, entities.OptionSetTypeArray, actual.StartingProficiencyOptions.OptionSet.GetType())
		assert.Equal(t, 3, len(actual.StartingProficiencyOptions.OptionSet.(*entities.OptionsArrayOptionSet).Options))
		assert.Equal(t, "smiths-tools", actual.StartingProficiencyOptions.OptionSet.(*entities.OptionsArrayOptionSet).Options[0].(*entities.ReferenceOption).Reference.Key)
		assert.Equal(t, "brewers-supplies", actual.StartingProficiencyOptions.OptionSet.(*entities.OptionsArrayOptionSet).Options[1].(*entities.ReferenceOption).Reference.Key)

	})
}

func TestDnd5eAPI_ListEquipment(t *testing.T) {
	t.Run("returns error if http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"equipment").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListEquipment()

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("returns error if decoding fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"equipment").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListEquipment()

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("returns error if status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"equipment").Return(&http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListEquipment()

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a list of equipment", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/equipment/equipmentlist.json")
		equipmentFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"equipment").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(equipmentFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		actual, err := dnd5eAPI.ListEquipment()

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, 237, len(actual))
		assert.Equal(t, "abacus", actual[0].Key)
		assert.Equal(t, "Abacus", actual[0].Name)
	})
}

func TestDnd5eAPI_GetEquipment(t *testing.T) {
	t.Run("returns error if http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"equipment/abacus").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetEquipment("abacus")

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("returns error if decoding fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"equipment/abacus").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetEquipment("abacus")

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("returns error if status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"equipment/abacus").Return(&http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetEquipment("abacus")

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns an equipment", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/equipment/abacus.json")
		equipmentFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"equipment/abacus").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(equipmentFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		actual, err := dnd5eAPI.GetEquipment("abacus")

		assert.Nil(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, "abacus", actual.Key)
		assert.Equal(t, "Abacus", actual.Name)
		assert.Equal(t, "adventuring-gear", actual.EquipmentCategory.Key)
		assert.Equal(t, "Adventuring Gear", actual.EquipmentCategory.Name)
		assert.Equal(t, 2, actual.Cost.Quantity)
		assert.Equal(t, "gp", actual.Cost.Unit)
		assert.Equal(t, 2, actual.Weight)
	})

}
