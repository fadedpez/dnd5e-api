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
	})

}
