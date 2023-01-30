package dnd5e

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
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
		assert.Equal(t, 1, actual.LanguageOptions.ChoiceCount)
		assert.Equal(t, "languages", actual.LanguageOptions.ChoiceType)
		assert.Equal(t, 15, len(actual.LanguageOptions.OptionList.Options))
		assert.Equal(t, entities.OptionTypeReference, actual.LanguageOptions.OptionList.Options[0].GetOptionType())
		assert.Equal(t, "dwarvish", actual.LanguageOptions.OptionList.Options[0].(*entities.ReferenceOption).Reference.Key)
		assert.Equal(t, entities.OptionTypeReference, actual.LanguageOptions.OptionList.Options[1].GetOptionType())
		assert.Equal(t, "elvish", actual.LanguageOptions.OptionList.Options[1].(*entities.ReferenceOption).Reference.Key)
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
		assert.Equal(t, 1, actual.StartingProficiencyOptions.ChoiceCount)
		assert.Equal(t, 3, len(actual.StartingProficiencyOptions.OptionList.Options))
		assert.Equal(t, entities.OptionTypeReference, actual.StartingProficiencyOptions.OptionList.Options[0].GetOptionType())
		assert.Equal(t, "smiths-tools", actual.StartingProficiencyOptions.OptionList.Options[0].(*entities.ReferenceOption).Reference.Key)
		assert.Equal(t, entities.OptionTypeReference, actual.StartingProficiencyOptions.OptionList.Options[1].GetOptionType())
		assert.Equal(t, "brewers-supplies", actual.StartingProficiencyOptions.OptionList.Options[1].(*entities.ReferenceOption).Reference.Key)

	})
}

func TestDND5eAPI_ListEquipment(t *testing.T) {
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

func TestDND5eAPI_GetEquipment(t *testing.T) {
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
		assert.Equal(t, "equipment", actual.GetType())
		equipment := actual.(*entities.Equipment)
		assert.Equal(t, "abacus", equipment.Key)
		assert.Equal(t, "Abacus", equipment.Name)
		assert.Equal(t, "adventuring-gear", equipment.EquipmentCategory.Key)
		assert.Equal(t, "Adventuring Gear", equipment.EquipmentCategory.Name)
		assert.Equal(t, 2, equipment.Cost.Quantity)
		assert.Equal(t, "gp", equipment.Cost.Unit)
		assert.Equal(t, 2, equipment.Weight)
	})

	t.Run("it returns a weapon", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/equipment/battleaxe.json")
		equipmentFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"equipment/battleaxe").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(equipmentFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetEquipment("battleaxe")

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "weapon", result.GetType())
		actual := result.(*entities.Weapon)
		assert.Equal(t, "battleaxe", actual.Key)
		assert.Equal(t, "Battleaxe", actual.Name)
		assert.Equal(t, "weapon", actual.EquipmentCategory.Key)
		assert.Equal(t, "Martial", actual.WeaponCategory)
		assert.Equal(t, "Melee", actual.WeaponRange)
		assert.Equal(t, "Martial Melee", actual.CategoryRange)
		assert.Equal(t, "1d8", actual.Damage.DamageDice)
		assert.Equal(t, "slashing", actual.Damage.DamageType.Key)
		assert.Equal(t, "Slashing", actual.Damage.DamageType.Name)
		assert.Equal(t, 5, actual.Range.Normal)
		assert.Equal(t, "versatile", actual.Properties[0].Key)
		assert.Equal(t, "Versatile", actual.Properties[0].Name)
		assert.Equal(t, "1d10", actual.TwoHandedDamage.DamageDice)
		assert.Equal(t, "slashing", actual.TwoHandedDamage.DamageType.Key)
		assert.Equal(t, "Slashing", actual.TwoHandedDamage.DamageType.Name)
		assert.Equal(t, 4, actual.Weight)
		assert.Equal(t, 10, actual.Cost.Quantity)
		assert.Equal(t, "gp", actual.Cost.Unit)
	})

	t.Run("it returns a weapon property", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/equipment/studded-leather-armor.json")
		equipmentFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"equipment/studded-leather-armor").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(equipmentFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetEquipment("studded-leather-armor")

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "armor", result.GetType())
		actual := result.(*entities.Armor)
		assert.Equal(t, "studded-leather-armor", actual.Key)
		assert.Equal(t, "Studded Leather Armor", actual.Name)
		assert.Equal(t, "armor", actual.EquipmentCategory.Key)
		assert.Equal(t, "Light", actual.ArmorCategory)
		assert.Equal(t, 12, actual.ArmorClass.Base)
		assert.Equal(t, true, actual.ArmorClass.DexBonus)
		assert.Equal(t, 0, actual.StrMinimum)
		assert.Equal(t, false, actual.StealthDisadvantage)
		assert.Equal(t, 13, actual.Weight)
		assert.Equal(t, 45, actual.Cost.Quantity)
		assert.Equal(t, "gp", actual.Cost.Unit)
	})
}

func TestDND5eAPI_ListClasses(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"classes").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListClasses()

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"classes").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListClasses()

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"classes").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListClasses()

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a list of classes", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/classlist.json")
		classesFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classesFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.ListClasses()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 12, len(result))
		assert.Equal(t, "barbarian", result[0].Key)
		assert.Equal(t, "Barbarian", result[0].Name)
	})
}

func TestDND5eAPI_GetClass(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"classes/ranger").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetClass("ranger")

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"classes/ranger").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetClass("ranger")

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"classes/ranger").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetClass("ranger")

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a class", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/ranger.json")
		classFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/ranger").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classFile)),
		}, nil)

		filePath, _ = filepath.Abs("../../testdata/equipment_categories/simplemeleeweapons.json")
		equipmentFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"equipment-categories/simple-melee-weapons").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(equipmentFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClass("ranger")

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "ranger", result.Key)
		assert.Equal(t, "Ranger", result.Name)
		assert.Equal(t, 10, result.HitDie)
		assert.Equal(t, "light-armor", result.Proficiencies[0].Key)
		assert.Equal(t, "medium-armor", result.Proficiencies[1].Key)
		assert.Equal(t, "str", result.SavingThrows[0].Key)
		assert.Equal(t, "STR", result.SavingThrows[0].Name)
		assert.Equal(t, 1, result.StartingEquipment[0].Quantity)
		assert.Equal(t, "longbow", result.StartingEquipment[0].Equipment.Key)
		assert.Equal(t, "Longbow", result.StartingEquipment[0].Equipment.Name)
		assert.Equal(t, 20, result.StartingEquipment[1].Quantity)
		assert.Equal(t, "arrow", result.StartingEquipment[1].Equipment.Key)
		assert.Equal(t, "Arrow", result.StartingEquipment[1].Equipment.Name)
		assert.Equal(t, 1, len(result.ProficiencyChoices))
		assert.Equal(t, 8, len(result.ProficiencyChoices[0].OptionList.Options))
		assert.Equal(t, 3, result.ProficiencyChoices[0].ChoiceCount)
		assert.Equal(t, "proficiencies", result.ProficiencyChoices[0].ChoiceType)
		assert.Equal(t, 3, len(result.StartingEquipmentOptions))
		assert.Equal(t, entities.OptionTypeChoice, result.StartingEquipmentOptions[1].OptionList.Options[1].GetOptionType())
		choiceOption := result.StartingEquipmentOptions[1].OptionList.Options[1].(*entities.ChoiceOption)
		assert.Equal(t, 10, len(choiceOption.OptionList.Options))
	})
}

func TestDnd5eAPI_ListSpells(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"spells").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListSpells(&ListSpellsInput{})

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"spells").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListSpells(&ListSpellsInput{})

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"spells").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListSpells(&ListSpellsInput{})

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a list of spells with no level or class specified", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/spells/spelllist.json")
		spellsFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"spells").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(spellsFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.ListSpells(&ListSpellsInput{})

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 319, len(result))
		assert.Equal(t, "acid-arrow", result[0].Key)
		assert.Equal(t, "Acid Arrow", result[0].Name)

	})

	t.Run("it returns a list of spells at level when level is specified and class is not", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/levels/1/spells.json")
		spellsFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"spells?level=1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(spellsFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		setLevel := 1
		result, err := dnd5eAPI.ListSpells(&ListSpellsInput{Level: &setLevel})

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 49, len(result))
		assert.Equal(t, "alarm", result[0].Key)
		assert.Equal(t, "Alarm", result[0].Name)
	})

	t.Run("it returns a list of spells of a class when level is not specified", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/sorcerer/sorcerer_spells.json")
		spellsFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/sorcerer/spells").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(spellsFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		setClass := "sorcerer"
		result, err := dnd5eAPI.ListSpells(&ListSpellsInput{Class: setClass})

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 120, len(result))
		assert.Equal(t, "acid-splash", result[0].Key)
		assert.Equal(t, "Acid Splash", result[0].Name)
	})

	t.Run("it returns a list of spells of a specified class and level", func(t *testing.T) {
		client := &mockHTTPClient{}
		classFilePath, _ := filepath.Abs("../../testdata/classes/sorcerer/sorcerer_spells.json")
		classSpellsFile, err := os.ReadFile(classFilePath)
		assert.Nil(t, err)

		levelFilePath, _ := filepath.Abs("../../testdata/levels/1/spells.json")
		levelSpellsFile, err := os.ReadFile(levelFilePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/sorcerer/spells").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classSpellsFile)),
		}, nil)

		client.On("Get", baserulzURL+"spells?level=1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(levelSpellsFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		setClass := "sorcerer"
		setLevel := 1
		result, err := dnd5eAPI.ListSpells(&ListSpellsInput{Class: setClass, Level: &setLevel})

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 17, len(result))
		assert.Equal(t, "burning-hands", result[0].Key)
		assert.Equal(t, "Burning Hands", result[0].Name)
		assert.Equal(t, "charm-person", result[1].Key)
		assert.Equal(t, "Charm Person", result[1].Name)
	})
}

func TestDND5eAPI_GetSpell(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"spells/burning-hands").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetSpell("burning-hands")

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"spells/burning-hands").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetSpell("burning-hands")

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"spells/burning-hands").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetSpell("burning-hands")

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a spell", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/spells/burninghands.json")
		spellFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"spells/burning-hands").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(spellFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetSpell("burning-hands")

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "burning-hands", result.Key)
		assert.Equal(t, "Burning Hands", result.Name)
		assert.Equal(t, "Self", result.Range)
		assert.Equal(t, false, result.Ritual)
		assert.Equal(t, "Instantaneous", result.Duration)
		assert.Equal(t, false, result.Concentration)
		assert.Equal(t, "1 action", result.CastingTime)
		assert.Equal(t, 1, result.SpellLevel)
		assert.Equal(t, "fire", result.SpellDamage.SpellDamageType.Key)
		assert.Equal(t, "Fire", result.SpellDamage.SpellDamageType.Name)
		assert.Equal(t, "9d6", result.SpellDamage.SpellDamageAtSlotLevel.SeventhLevel)
		assert.Equal(t, "dex", result.DC.DCType.Key)
		assert.Equal(t, "DEX", result.DC.DCType.Name)
		assert.Equal(t, "half", result.DC.DCSuccess)
		assert.Equal(t, "cone", result.AreaOfEffect.Type)
		assert.Equal(t, 15, result.AreaOfEffect.Size)
		assert.Equal(t, "evocation", result.SpellSchool.Key)
		assert.Equal(t, "Evocation", result.SpellSchool.Name)
		assert.Equal(t, "sorcerer", result.SpellClasses[0].Key)
		assert.Equal(t, "wizard", result.SpellClasses[1].Key)
	})
}

func TestDND5eAPI_ListFeatures(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"features").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListFeatures()

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"features").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListFeatures()

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"features").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListFeatures()

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a list of features", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/features/featurelist.json")
		featuresFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"features").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(featuresFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.ListFeatures()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 370, len(result))
		assert.Equal(t, "action-surge-1-use", result[0].Key)
		assert.Equal(t, "Action Surge (1 use)", result[0].Name)
	})
}

func TestDND5eAPI_GetFeature(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"features/metamagic-2").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetFeature("metamagic-2")

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"features/metamagic-2").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetFeature("metamagic-2")

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"features/metamagic-2").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetFeature("metamagic-2")

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a feature", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/features/metamagic2.json")
		featureFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"features/metamagic-2").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(featureFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetFeature("metamagic-2")

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "metamagic-2", result.Key)
		assert.Equal(t, "Metamagic", result.Name)
		assert.Equal(t, 10, result.Level)
		assert.Equal(t, "sorcerer", result.Class.Key)
		assert.Equal(t, "Sorcerer", result.Class.Name)
		assert.Equal(t, 1, result.FeatureSpecific.SubFeatureOptions.ChoiceCount)
		assert.Equal(t, entities.OptionTypeReference, result.FeatureSpecific.SubFeatureOptions.OptionList.Options[0].GetOptionType())
		assert.Equal(t, 8, len(result.FeatureSpecific.SubFeatureOptions.OptionList.Options))
		assert.Equal(t, "metamagic-careful-spell", result.FeatureSpecific.SubFeatureOptions.OptionList.Options[0].(*entities.ReferenceOption).Reference.Key)
		assert.Equal(t, "Metamagic: Careful Spell", result.FeatureSpecific.SubFeatureOptions.OptionList.Options[0].(*entities.ReferenceOption).Reference.Name)
		assert.Equal(t, entities.OptionTypeReference, result.FeatureSpecific.SubFeatureOptions.OptionList.Options[1].GetOptionType())
		assert.Equal(t, "metamagic-distant-spell", result.FeatureSpecific.SubFeatureOptions.OptionList.Options[1].(*entities.ReferenceOption).Reference.Key)
		assert.Equal(t, "Metamagic: Distant Spell", result.FeatureSpecific.SubFeatureOptions.OptionList.Options[1].(*entities.ReferenceOption).Reference.Name)
	})
}

func TestDND5eAPI_ListSkills(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"skills").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListSkills()

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"skills").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListSkills()

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"skills").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListSkills()

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a list of skills", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/skills/skilllist.json")
		skillsFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"skills").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(skillsFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.ListSkills()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 18, len(result))
		assert.Equal(t, "acrobatics", result[0].Key)
		assert.Equal(t, "Acrobatics", result[0].Name)
	})
}

func TestDND5eAPI_GetSkill(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"skills/acrobatics").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetSkill("acrobatics")

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"skills/acrobatics").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetSkill("acrobatics")

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"skills/acrobatics").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetSkill("acrobatics")

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a skill", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/skills/acrobatics.json")
		skillFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"skills/acrobatics").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(skillFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetSkill("acrobatics")

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "acrobatics", result.Key)
		assert.Equal(t, "Acrobatics", result.Name)
		assert.Equal(t, "dex", result.AbilityScore.Key)
		assert.Equal(t, "DEX", result.AbilityScore.Name)
		assert.Equal(t, "skills", result.Type)
	})
}

func TestDND5eAPI_ListMonsters(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"monsters").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListMonsters()

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"monsters").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListMonsters()

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"monsters").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListMonsters()

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a list of monsters", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/monsters/monsterslist.json")
		monstersFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"monsters").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(monstersFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.ListMonsters()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 334, len(result))
		assert.Equal(t, "aboleth", result[0].Key)
		assert.Equal(t, "Aboleth", result[0].Name)
	})
}

func TestDND5eAPI_GetMonster(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"monsters/goblin").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetMonster("goblin")

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"monsters/goblin").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetMonster("goblin")

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"monsters/goblin").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetMonster("goblin")

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a monster", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/monsters/goblin.json")
		monsterFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"monsters/goblin").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(monsterFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetMonster("goblin")

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "goblin", result.Key)
		assert.Equal(t, "Goblin", result.Name)
		assert.Equal(t, "Small", result.Size)
		assert.Equal(t, "humanoid", result.Type)
		assert.Equal(t, "neutral evil", result.Alignment)
		assert.Equal(t, 15, result.ArmorClass)
		assert.Equal(t, 7, result.HitPoints)
		assert.Equal(t, "2d6", result.HitDice)
		assert.Equal(t, "30 ft.", result.Speed.Walk)
		assert.Equal(t, 8, result.Strength)
		assert.Equal(t, 14, result.Dexterity)
		assert.Equal(t, 10, result.Constitution)
		assert.Equal(t, 10, result.Intelligence)
		assert.Equal(t, 8, result.Wisdom)
		assert.Equal(t, 8, result.Charisma)
		assert.Equal(t, 6, result.Proficiencies[0].Value)
		assert.Equal(t, "skill-stealth", result.Proficiencies[0].Proficiency.Key)
		assert.Equal(t, "Skill: Stealth", result.Proficiencies[0].Proficiency.Name)
		assert.Equal(t, "60 ft.", result.MonsterSenses.Darkvision)
		assert.Equal(t, 9, result.MonsterSenses.PassivePerception)
		assert.Equal(t, "Common, Goblin", result.Languages)
		assert.Equal(t, float32(0.25), result.ChallengeRating)
		assert.Equal(t, 50, result.XP)
		assert.Equal(t, "Scimitar", result.MonsterActions[0].Name)
		assert.Equal(t, "Melee Weapon Attack: +4 to hit, reach 5 ft., one target. Hit: 5 (1d6 + 2) slashing damage.", result.MonsterActions[0].Description)
		assert.Equal(t, 4, result.MonsterActions[0].AttackBonus)
		assert.Equal(t, "slashing", result.MonsterActions[0].Damage[0].DamageType.Key)
		assert.Equal(t, "Slashing", result.MonsterActions[0].Damage[0].DamageType.Name)
		assert.Equal(t, "1d6+2", result.MonsterActions[0].Damage[0].DamageDice)
		assert.Equal(t, "/api/images/monsters/goblin.png", result.MonsterImageURL)
		assert.Equal(t, "sandwiches", result.DamageVulnerabilities[0])
		assert.Equal(t, "lightning", result.DamageResistances[0])
		assert.Equal(t, "fire", result.DamageImmunities[0])
		assert.Equal(t, "exhaustion", result.ConditionImmunities[0].Key)
		assert.Equal(t, "Exhaustion", result.ConditionImmunities[0].Name)
	})
}

func TestDND5eAPI_GetClassLevel(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"classes/ranger/levels/1").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetClassLevel("ranger", 1)

		assert.NotNil(t, err)
		assert.Equal(t, "http.Get failed", err.Error())
	})

	t.Run("it returns an error if json.Unmarshal fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"classes/ranger/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetClassLevel("ranger", 1)

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"classes/ranger/levels/1").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetClassLevel("ranger", 1)

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a class level", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/rangerlevel1.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/ranger/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("ranger", 1)

		assert.Nil(t, err)
		assert.Equal(t, 1, result.Level)
		assert.Equal(t, 2, result.ProfBonus)
		assert.Equal(t, 2, len(result.Features))
		assert.Equal(t, "Favored Enemy (1 type)", result.Features[0].Name)
		assert.Equal(t, 0, result.SpellCasting.SpellsKnown)
		assert.Equal(t, 0, result.SpellCasting.SpellSlotsLevel1)
		assert.Equal(t, 0, result.SpellCasting.SpellSlotsLevel5)
		assert.Equal(t, "ranger-1", result.Key)
		assert.Equal(t, "ranger", result.Class.Key)
		assert.Equal(t, "Ranger", result.Class.Name)
		assert.Equal(t, "ranger", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 1, result.ClassSpecific.(*entities.RangerSpecific).FavoredEnemies)
	})

	t.Run("it returns barbarian specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/barbarianlevel1.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/barbarian/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("barbarian", 1)

		assert.Nil(t, err)
		assert.Equal(t, 1, result.Level)
		assert.Equal(t, "barbarian", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 2, result.ClassSpecific.(*entities.BarbarianSpecific).RageCount)
		assert.Equal(t, 2, result.ClassSpecific.(*entities.BarbarianSpecific).RageDamageBonus)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.BarbarianSpecific).BrutalCriticalDice)
	})

	t.Run("it returns bard specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/bardlevel1.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/bard/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("bard", 1)

		assert.Nil(t, err)
		assert.Equal(t, "bard", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 6, result.ClassSpecific.(*entities.BardSpecific).BardicInspirationDie)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.BardSpecific).SongOfRestDie)
	})

	t.Run("it returns cleric specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/clericlevel1.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/cleric/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("cleric", 1)

		assert.Nil(t, err)
		assert.Equal(t, 3, result.SpellCasting.CantripsKnown)
		assert.Equal(t, "cleric", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 0, result.ClassSpecific.(*entities.ClericSpecific).ChannelDivinityCharges)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.ClericSpecific).DestroyUndeadCR)
	})

	t.Run("it returns druid specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/druidlevel1.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/druid/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("druid", 1)

		assert.Nil(t, err)
		assert.Equal(t, "druid", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 0, result.ClassSpecific.(*entities.DruidSpecific).WildShapeMaxCR)
		assert.Equal(t, false, result.ClassSpecific.(*entities.DruidSpecific).WildShapeSwim)
		assert.Equal(t, false, result.ClassSpecific.(*entities.DruidSpecific).WildShapeFly)
	})

	t.Run("it returns fighter specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/fighterlevel1.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/fighter/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("fighter", 1)

		assert.Nil(t, err)
		assert.Equal(t, "fighter", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 0, result.ClassSpecific.(*entities.FighterSpecific).ActionSurges)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.FighterSpecific).IndomitableUses)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.FighterSpecific).ExtraAttacks)
	})

	t.Run("it returns monk specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/monklevel1.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/monk/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("monk", 1)

		assert.Nil(t, err)
		assert.Equal(t, "monk", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 1, result.ClassSpecific.(*entities.MonkSpecific).MartialArts.DiceCount)
		assert.Equal(t, 4, result.ClassSpecific.(*entities.MonkSpecific).MartialArts.DiceValue)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.MonkSpecific).KiPoints)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.MonkSpecific).UnarmoredMovement)
	})

	t.Run("it returns paladin specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/paladinlevel1.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/paladin/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("paladin", 1)

		assert.Nil(t, err)
		assert.Equal(t, "paladin", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 0, result.ClassSpecific.(*entities.PaladinSpecific).AuraRange)
	})

	t.Run("it returns rogue specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/roguelevel1.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/rogue/levels/1").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("rogue", 1)

		assert.Nil(t, err)
		assert.Equal(t, "rogue", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 1, result.ClassSpecific.(*entities.RogueSpecific).SneakAttack.DiceCount)
		assert.Equal(t, 6, result.ClassSpecific.(*entities.RogueSpecific).SneakAttack.DiceValue)
	})

	t.Run("it returns sorcerer specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/sorcererlevel5.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/sorcerer/levels/5").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("sorcerer", 5)

		assert.Nil(t, err)
		assert.Equal(t, "sorcerer", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 5, result.ClassSpecific.(*entities.SorcererSpecific).SorceryPoints)
		assert.Equal(t, 2, result.ClassSpecific.(*entities.SorcererSpecific).MetamagicKnown)
		assert.Equal(t, 1, result.ClassSpecific.(*entities.SorcererSpecific).CreatingSpellSlots[0].SpellSlotLevel)
		assert.Equal(t, 2, result.ClassSpecific.(*entities.SorcererSpecific).CreatingSpellSlots[0].SorceryPointCost)
	})

	t.Run("it returns warlock specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/warlocklevel5.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/warlock/levels/5").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("warlock", 5)

		assert.Nil(t, err)
		assert.Equal(t, "warlock", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 3, result.ClassSpecific.(*entities.WarlockSpecific).InvocationsKnown)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.WarlockSpecific).MysticArcanumLevel6)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.WarlockSpecific).MysticArcanumLevel7)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.WarlockSpecific).MysticArcanumLevel8)
		assert.Equal(t, 0, result.ClassSpecific.(*entities.WarlockSpecific).MysticArcanumLevel9)
	})

	t.Run("it returns wizard specific level data", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/classes/levels/wizardlevel5.json")
		classLevelFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"classes/wizard/levels/5").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(classLevelFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.GetClassLevel("wizard", 5)

		assert.Nil(t, err)
		assert.Equal(t, "wizard", result.ClassSpecific.GetSpecificClass())
		assert.Equal(t, 3, result.ClassSpecific.(*entities.WizardSpecific).ArcaneRecoveryLevels)
	})
}

func TestDnd5eAPI_GetProficiency(t *testing.T) {
	type fields struct {
		client *mockHTTPClient
	}

	type args struct {
		key string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entities.Proficiency
		wantErr bool
	}{
		{
			name: "it returns a skill proficiency",
			fields: fields{
				client: &mockHTTPClient{},
			},
			args: args{
				key: "skill-animal-handling",
			},
			want: &entities.Proficiency{
				Key:  "skill-animal-handling",
				Name: "Skill: Animal Handling",
				Type: entities.ProficiencyTypeSkill,
				Reference: &entities.ReferenceItem{
					Key:  "animal-handling",
					Name: "Animal Handling",
					Type: "skills",
				},
			},
			wantErr: false,
		}, {
			name: "it returns a tool proficiency",
			fields: fields{
				client: &mockHTTPClient{},
			},
			args: args{
				key: "smiths-tools",
			},
			want: &entities.Proficiency{
				Key:  "smiths-tools",
				Name: "Smith's Tools",
				Type: entities.ProficiencyTypeTool,
				Reference: &entities.ReferenceItem{
					Key:  "smiths-tools",
					Name: "Smith's Tools",
					Type: "equipment",
				},
			},
			wantErr: false,
		}, {
			name: "it returns a saving throw proficiency",
			fields: fields{
				client: &mockHTTPClient{},
			},
			args: args{
				key: "saving-throw-str",
			},
			want: &entities.Proficiency{
				Key:  "saving-throw-str",
				Name: "Saving Throw: STR",
				Type: entities.ProficiencyTypeSavingThrow,
				Reference: &entities.ReferenceItem{
					Key:  "str",
					Name: "STR",
					Type: "ability-scores",
				},
			},
			wantErr: false,
		}, {
			name: "it returns an armor proficiency",
			fields: fields{
				client: &mockHTTPClient{},
			},
			args: args{
				key: "light-armor",
			},
			want: &entities.Proficiency{
				Key:  "light-armor",
				Name: "Light Armor",
				Type: entities.ProficiencyTypeArmor,
				Reference: &entities.ReferenceItem{
					Key:  "light-armor",
					Name: "Light Armor",
					Type: "equipment-categories",
				},
			},
			wantErr: false,
		}, {
			name: "it returns a weapon proficiency",
			fields: fields{
				client: &mockHTTPClient{},
			},
			args: args{
				key: "simple-weapons",
			},
			want: &entities.Proficiency{
				Key:  "simple-weapons",
				Name: "Simple Weapons",
				Type: entities.ProficiencyTypeWeapon,
				Reference: &entities.ReferenceItem{
					Key:  "simple-weapons",
					Name: "Simple Weapons",
					Type: "equipment-categories",
				},
			},
		}, {
			name: "it returns an instrument proficiency",
			fields: fields{
				client: &mockHTTPClient{},
			},
			args: args{
				key: "lute",
			},
			want: &entities.Proficiency{
				Key:  "lute",
				Name: "Lute",
				Type: entities.ProficiencyTypeInstrument,
				Reference: &entities.ReferenceItem{
					Key:  "lute",
					Name: "Lute",
					Type: "equipment",
				},
			},
		}, {
			name: "it returns an error on invalid json",
			fields: fields{
				client: &mockHTTPClient{},
			},
			args: args{
				key: "invalid",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath, _ := filepath.Abs(fmt.Sprintf("../../testdata/proficiencies/%s.json", tt.args.key))
			proficiencyFile, err := os.ReadFile(filePath)
			tt.fields.client.On("Get", baserulzURL+"proficiencies/"+tt.args.key).Return(&http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewReader(proficiencyFile)),
			}, nil)

			dnd5eAPI := &dnd5eAPI{
				client: tt.fields.client,
			}

			got, err := dnd5eAPI.GetProficiency(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProficiency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProficiency() got = %v, want %v", got, tt.want)
			}
		})
	}
}
