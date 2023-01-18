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
		assert.Equal(t, 15, len(actual.LanguageOptions.OptionSet.Options))
		assert.Equal(t, "dwarvish", actual.LanguageOptions.OptionSet.Options[0].(*entities.ReferenceOption).Reference.Key)
		assert.Equal(t, "elvish", actual.LanguageOptions.OptionSet.Options[1].(*entities.ReferenceOption).Reference.Key)
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
		assert.Equal(t, 3, len(actual.StartingProficiencyOptions.OptionSet.Options))
		assert.Equal(t, "smiths-tools", actual.StartingProficiencyOptions.OptionSet.Options[0].(*entities.ReferenceOption).Reference.Key)
		assert.Equal(t, "brewers-supplies", actual.StartingProficiencyOptions.OptionSet.Options[1].(*entities.ReferenceOption).Reference.Key)

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
		assert.Equal(t, 8, len(result.ProficiencyChoices[0].OptionSet.Options))
		assert.Equal(t, 3, result.ProficiencyChoices[0].Choose)
		assert.Equal(t, "proficiencies", result.ProficiencyChoices[0].Type)
	})
}

func TestDnd5eAPI_ListSpells(t *testing.T) {
	t.Run("it returns an error when http.Get fails", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"spells").Return(nil, errors.New("http.Get failed"))

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListSpells()

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
		_, err := dnd5eAPI.ListSpells()

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"spells").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.ListSpells()

		assert.NotNil(t, err)
		assert.Equal(t, "unexpected status code: 500", err.Error())
	})

	t.Run("it returns a list of spells", func(t *testing.T) {
		client := &mockHTTPClient{}
		filePath, _ := filepath.Abs("../../testdata/spells/spelllist.json")
		spellsFile, err := os.ReadFile(filePath)
		assert.Nil(t, err)

		client.On("Get", baserulzURL+"spells").Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(spellsFile)),
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		result, err := dnd5eAPI.ListSpells()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 319, len(result))
		assert.Equal(t, "acid-arrow", result[0].Key)
		assert.Equal(t, "Acid Arrow", result[0].Name)

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
		_, err := dnd5eAPI.GetFeatures()

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
		_, err := dnd5eAPI.GetFeatures()

		assert.NotNil(t, err)
		assert.Equal(t, "invalid character 'i' looking for beginning of value", err.Error())
	})

	t.Run("it returns an error when the status code is not 200", func(t *testing.T) {
		client := &mockHTTPClient{}
		client.On("Get", baserulzURL+"features").Return(&http.Response{
			StatusCode: 500,
		}, nil)

		dnd5eAPI := &dnd5eAPI{client: client}
		_, err := dnd5eAPI.GetFeatures()

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
		result, err := dnd5eAPI.GetFeatures()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, 370, len(result))
		assert.Equal(t, "action-surge-1-use", result[0].Key)
		assert.Equal(t, "Action Surge (1 use)", result[0].Name)
	})
}
