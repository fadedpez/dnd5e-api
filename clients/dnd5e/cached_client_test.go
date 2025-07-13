package dnd5e

import (
	"errors"
	"testing"
	"time"

	"github.com/fadedpez/dnd5e-api/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockClient is a mock implementation of the Interface
type MockClient struct {
	mock.Mock
}

func (m *MockClient) ListRaces() ([]*entities.ReferenceItem, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReferenceItem), args.Error(1)
}

func (m *MockClient) GetRace(key string) (*entities.Race, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Race), args.Error(1)
}

func (m *MockClient) ListEquipment() ([]*entities.ReferenceItem, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReferenceItem), args.Error(1)
}

func (m *MockClient) GetEquipment(key string) (EquipmentInterface, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(EquipmentInterface), args.Error(1)
}

func (m *MockClient) ListClasses() ([]*entities.ReferenceItem, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReferenceItem), args.Error(1)
}

func (m *MockClient) GetClass(key string) (*entities.Class, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Class), args.Error(1)
}

func (m *MockClient) ListSpells(input *ListSpellsInput) ([]*entities.ReferenceItem, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReferenceItem), args.Error(1)
}

func (m *MockClient) GetSpell(key string) (*entities.Spell, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Spell), args.Error(1)
}

func (m *MockClient) ListFeatures() ([]*entities.ReferenceItem, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReferenceItem), args.Error(1)
}

func (m *MockClient) GetFeature(key string) (*entities.Feature, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Feature), args.Error(1)
}

func (m *MockClient) ListSkills() ([]*entities.ReferenceItem, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReferenceItem), args.Error(1)
}

func (m *MockClient) GetSkill(key string) (*entities.Skill, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Skill), args.Error(1)
}

func (m *MockClient) ListMonsters() ([]*entities.ReferenceItem, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReferenceItem), args.Error(1)
}

func (m *MockClient) ListMonstersWithFilter(input *ListMonstersInput) ([]*entities.ReferenceItem, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReferenceItem), args.Error(1)
}

func (m *MockClient) GetMonster(key string) (*entities.Monster, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Monster), args.Error(1)
}

func (m *MockClient) GetClassLevel(key string, level int) (*entities.Level, error) {
	args := m.Called(key, level)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Level), args.Error(1)
}

func (m *MockClient) GetProficiency(key string) (*entities.Proficiency, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Proficiency), args.Error(1)
}

func (m *MockClient) ListDamageTypes() ([]*entities.ReferenceItem, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.ReferenceItem), args.Error(1)
}

func (m *MockClient) GetDamageType(key string) (*entities.DamageType, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.DamageType), args.Error(1)
}

func (m *MockClient) GetEquipmentCategory(key string) (*entities.EquipmentCategory, error) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.EquipmentCategory), args.Error(1)
}

func TestCachedClient_GetRace_CacheHit(t *testing.T) {
	mockClient := new(MockClient)
	cachedClient := NewCachedClient(mockClient, 24*time.Hour).(*CachedClient)

	// Setup expected race
	expectedRace := &entities.Race{
		Key:  "dwarf",
		Name: "Dwarf",
	}

	// First call - should hit the API
	mockClient.On("GetRace", "dwarf").Return(expectedRace, nil).Once()
	
	race1, err1 := cachedClient.GetRace("dwarf")
	assert.NoError(t, err1)
	assert.Equal(t, expectedRace, race1)

	// Second call - should hit the cache (no new mock expectation)
	race2, err2 := cachedClient.GetRace("dwarf")
	assert.NoError(t, err2)
	assert.Equal(t, expectedRace, race2)

	// Verify the API was only called once
	mockClient.AssertExpectations(t)
}

func TestCachedClient_GetRace_CacheMiss(t *testing.T) {
	mockClient := new(MockClient)
	cachedClient := NewCachedClient(mockClient, 24*time.Hour).(*CachedClient)

	// Setup expected races
	dwarf := &entities.Race{Key: "dwarf", Name: "Dwarf"}
	elf := &entities.Race{Key: "elf", Name: "Elf"}

	// Each different key should hit the API
	mockClient.On("GetRace", "dwarf").Return(dwarf, nil).Once()
	mockClient.On("GetRace", "elf").Return(elf, nil).Once()
	
	race1, err1 := cachedClient.GetRace("dwarf")
	assert.NoError(t, err1)
	assert.Equal(t, dwarf, race1)

	race2, err2 := cachedClient.GetRace("elf")
	assert.NoError(t, err2)
	assert.Equal(t, elf, race2)

	mockClient.AssertExpectations(t)
}

func TestCachedClient_GetRace_Error(t *testing.T) {
	mockClient := new(MockClient)
	cachedClient := NewCachedClient(mockClient, 24*time.Hour).(*CachedClient)

	expectedError := errors.New("API error")
	mockClient.On("GetRace", "invalid").Return(nil, expectedError).Once()
	
	race, err := cachedClient.GetRace("invalid")
	assert.Error(t, err)
	assert.Nil(t, race)
	assert.Equal(t, expectedError, err)

	// Error responses should not be cached
	mockClient.On("GetRace", "invalid").Return(nil, expectedError).Once()
	race2, err2 := cachedClient.GetRace("invalid")
	assert.Error(t, err2)
	assert.Nil(t, race2)

	mockClient.AssertExpectations(t)
}

func TestCachedClient_ListSpells_DifferentFilters(t *testing.T) {
	mockClient := new(MockClient)
	cachedClient := NewCachedClient(mockClient, 24*time.Hour).(*CachedClient)

	// Different spell lists for different filters
	allSpells := []*entities.ReferenceItem{{Key: "magic-missile", Name: "Magic Missile"}}
	wizardSpells := []*entities.ReferenceItem{{Key: "fireball", Name: "Fireball"}}
	level1Spells := []*entities.ReferenceItem{{Key: "shield", Name: "Shield"}}

	// Test nil input (all spells)
	mockClient.On("ListSpells", (*ListSpellsInput)(nil)).Return(allSpells, nil).Once()
	result1, err1 := cachedClient.ListSpells(nil)
	assert.NoError(t, err1)
	assert.Equal(t, allSpells, result1)

	// Test class filter
	wizardInput := &ListSpellsInput{Class: "wizard"}
	mockClient.On("ListSpells", wizardInput).Return(wizardSpells, nil).Once()
	result2, err2 := cachedClient.ListSpells(wizardInput)
	assert.NoError(t, err2)
	assert.Equal(t, wizardSpells, result2)

	// Test level filter
	level := 1
	levelInput := &ListSpellsInput{Level: &level}
	mockClient.On("ListSpells", levelInput).Return(level1Spells, nil).Once()
	result3, err3 := cachedClient.ListSpells(levelInput)
	assert.NoError(t, err3)
	assert.Equal(t, level1Spells, result3)

	// Verify each was cached separately
	// Call again with same filters - should not hit API
	result4, err4 := cachedClient.ListSpells(nil)
	assert.NoError(t, err4)
	assert.Equal(t, allSpells, result4)

	result5, err5 := cachedClient.ListSpells(wizardInput)
	assert.NoError(t, err5)
	assert.Equal(t, wizardSpells, result5)

	mockClient.AssertExpectations(t)
}

func TestCachedClient_CacheExpiration(t *testing.T) {
	mockClient := new(MockClient)
	// Use a very short TTL for testing
	cachedClient := NewCachedClient(mockClient, 100*time.Millisecond).(*CachedClient)

	expectedRace := &entities.Race{Key: "dwarf", Name: "Dwarf"}

	// First call - should hit the API
	mockClient.On("GetRace", "dwarf").Return(expectedRace, nil).Once()
	race1, err1 := cachedClient.GetRace("dwarf")
	assert.NoError(t, err1)
	assert.Equal(t, expectedRace, race1)

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Second call - should hit the API again due to expiration
	mockClient.On("GetRace", "dwarf").Return(expectedRace, nil).Once()
	race2, err2 := cachedClient.GetRace("dwarf")
	assert.NoError(t, err2)
	assert.Equal(t, expectedRace, race2)

	mockClient.AssertExpectations(t)
}

func TestCachedClient_GetMonster_CacheHit(t *testing.T) {
	mockClient := new(MockClient)
	cachedClient := NewCachedClient(mockClient, 24*time.Hour).(*CachedClient)

	expectedMonster := &entities.Monster{
		Key:  "goblin",
		Name: "Goblin",
	}

	// First call - should hit the API
	mockClient.On("GetMonster", "goblin").Return(expectedMonster, nil).Once()
	
	monster1, err1 := cachedClient.GetMonster("goblin")
	assert.NoError(t, err1)
	assert.Equal(t, expectedMonster, monster1)

	// Second call - should hit the cache
	monster2, err2 := cachedClient.GetMonster("goblin")
	assert.NoError(t, err2)
	assert.Equal(t, expectedMonster, monster2)

	mockClient.AssertExpectations(t)
}

func TestCachedClient_ListMonstersWithFilter(t *testing.T) {
	mockClient := new(MockClient)
	cachedClient := NewCachedClient(mockClient, 24*time.Hour).(*CachedClient)

	// Test different CR filters
	cr1 := 1.0
	cr2 := 2.0
	
	monsters1 := []*entities.ReferenceItem{{Key: "goblin", Name: "Goblin"}}
	monsters2 := []*entities.ReferenceItem{{Key: "ogre", Name: "Ogre"}}

	input1 := &ListMonstersInput{ChallengeRating: &cr1}
	input2 := &ListMonstersInput{ChallengeRating: &cr2}

	mockClient.On("ListMonstersWithFilter", input1).Return(monsters1, nil).Once()
	mockClient.On("ListMonstersWithFilter", input2).Return(monsters2, nil).Once()

	// Different filters should have separate cache entries
	result1, err1 := cachedClient.ListMonstersWithFilter(input1)
	assert.NoError(t, err1)
	assert.Equal(t, monsters1, result1)

	result2, err2 := cachedClient.ListMonstersWithFilter(input2)
	assert.NoError(t, err2)
	assert.Equal(t, monsters2, result2)

	// Calling again with same filter should use cache
	result3, err3 := cachedClient.ListMonstersWithFilter(input1)
	assert.NoError(t, err3)
	assert.Equal(t, monsters1, result3)

	mockClient.AssertExpectations(t)
}