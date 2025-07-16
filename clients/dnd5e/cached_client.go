package dnd5e

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/fadedpez/dnd5e-api/entities"
)

// cacheEntry holds cached data with timestamp
type cacheEntry struct {
	data      interface{}
	timestamp time.Time
}

// CachedClient wraps the D&D 5e API client with an in-memory cache
type CachedClient struct {
	client Interface
	cache  sync.Map
	ttl    time.Duration
}

// NewCachedClient creates a new cached client with specified TTL
func NewCachedClient(client Interface, ttl time.Duration) Interface {
	return &CachedClient{
		client: client,
		ttl:    ttl,
	}
}

// isExpired checks if a cache entry has expired
func (e *cacheEntry) isExpired(ttl time.Duration) bool {
	// For D&D 5e static data, we could use a very long TTL or no expiration
	// But keeping TTL for flexibility (e.g., 24 hours)
	return time.Since(e.timestamp) > ttl
}

// getFromCache attempts to retrieve and type-assert cached data
func (c *CachedClient) getFromCache(key string) (interface{}, bool) {
	if cached, ok := c.cache.Load(key); ok {
		if entry, ok := cached.(*cacheEntry); ok && !entry.isExpired(c.ttl) {
			return entry.data, true
		}
		// Remove expired entry
		c.cache.Delete(key)
	}
	return nil, false
}

// storeInCache stores data in the cache
func (c *CachedClient) storeInCache(key string, data interface{}) {
	c.cache.Store(key, &cacheEntry{
		data:      data,
		timestamp: time.Now(),
	})
}

// ListRaces returns cached race list or fetches from API
func (c *CachedClient) ListRaces() ([]*entities.ReferenceItem, error) {
	cacheKey := "list:races"
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	races, err := c.client.ListRaces()
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, races)
	return races, nil
}

// GetRace returns cached race or fetches from API
func (c *CachedClient) GetRace(key string) (*entities.Race, error) {
	cacheKey := fmt.Sprintf("race:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.Race); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.Race, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	race, err := c.client.GetRace(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, race)
	return race, nil
}

// ListEquipment returns cached equipment list or fetches from API
func (c *CachedClient) ListEquipment() ([]*entities.ReferenceItem, error) {
	cacheKey := "list:equipment"
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	equipment, err := c.client.ListEquipment()
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, equipment)
	return equipment, nil
}

// GetEquipment returns cached equipment or fetches from API
func (c *CachedClient) GetEquipment(key string) (EquipmentInterface, error) {
	cacheKey := fmt.Sprintf("equipment:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(EquipmentInterface); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected EquipmentInterface, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	equipment, err := c.client.GetEquipment(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, equipment)
	return equipment, nil
}

// ListClasses returns cached class list or fetches from API
func (c *CachedClient) ListClasses() ([]*entities.ReferenceItem, error) {
	cacheKey := "list:classes"
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	classes, err := c.client.ListClasses()
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, classes)
	return classes, nil
}

// GetClass returns cached class or fetches from API
func (c *CachedClient) GetClass(key string) (*entities.Class, error) {
	cacheKey := fmt.Sprintf("class:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.Class); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.Class, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	class, err := c.client.GetClass(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, class)
	return class, nil
}

// ListSpells returns cached spell list or fetches from API
func (c *CachedClient) ListSpells(input *ListSpellsInput) ([]*entities.ReferenceItem, error) {
	// Create unique cache key based on input parameters
	var cacheKey string
	if input == nil {
		cacheKey = "list:spells:all"
	} else if input.Class == "" && input.Level == nil {
		cacheKey = "list:spells:all"
	} else if input.Class == "" {
		cacheKey = fmt.Sprintf("list:spells:level:%d", *input.Level)
	} else if input.Level == nil {
		cacheKey = fmt.Sprintf("list:spells:class:%s", input.Class)
	} else {
		cacheKey = fmt.Sprintf("list:spells:class:%s:level:%d", input.Class, *input.Level)
	}
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	spells, err := c.client.ListSpells(input)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, spells)
	return spells, nil
}

// GetSpell returns cached spell or fetches from API
func (c *CachedClient) GetSpell(key string) (*entities.Spell, error) {
	cacheKey := fmt.Sprintf("spell:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.Spell); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.Spell, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	spell, err := c.client.GetSpell(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, spell)
	return spell, nil
}

// ListFeatures returns cached feature list or fetches from API
func (c *CachedClient) ListFeatures() ([]*entities.ReferenceItem, error) {
	cacheKey := "list:features"
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	features, err := c.client.ListFeatures()
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, features)
	return features, nil
}

// GetFeature returns cached feature or fetches from API
func (c *CachedClient) GetFeature(key string) (*entities.Feature, error) {
	cacheKey := fmt.Sprintf("feature:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.Feature); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.Feature, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	feature, err := c.client.GetFeature(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, feature)
	return feature, nil
}

// ListSkills returns cached skill list or fetches from API
func (c *CachedClient) ListSkills() ([]*entities.ReferenceItem, error) {
	cacheKey := "list:skills"
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	skills, err := c.client.ListSkills()
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, skills)
	return skills, nil
}

// GetSkill returns cached skill or fetches from API
func (c *CachedClient) GetSkill(key string) (*entities.Skill, error) {
	cacheKey := fmt.Sprintf("skill:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.Skill); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.Skill, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	skill, err := c.client.GetSkill(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, skill)
	return skill, nil
}

// ListMonsters returns cached monster list or fetches from API
func (c *CachedClient) ListMonsters() ([]*entities.ReferenceItem, error) {
	cacheKey := "list:monsters:all"
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	monsters, err := c.client.ListMonsters()
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, monsters)
	return monsters, nil
}

// ListMonstersWithFilter returns cached filtered monster list or fetches from API
func (c *CachedClient) ListMonstersWithFilter(input *ListMonstersInput) ([]*entities.ReferenceItem, error) {
	var cacheKey string
	if input == nil || input.ChallengeRating == nil {
		return c.ListMonsters()
	}
	
	cacheKey = fmt.Sprintf("list:monsters:cr:%g", *input.ChallengeRating)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	monsters, err := c.client.ListMonstersWithFilter(input)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, monsters)
	return monsters, nil
}

// GetMonster returns cached monster or fetches from API
func (c *CachedClient) GetMonster(key string) (*entities.Monster, error) {
	cacheKey := fmt.Sprintf("monster:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.Monster); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.Monster, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	monster, err := c.client.GetMonster(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, monster)
	return monster, nil
}

// GetClassLevel returns cached class level or fetches from API
func (c *CachedClient) GetClassLevel(key string, level int) (*entities.Level, error) {
	cacheKey := fmt.Sprintf("class:%s:level:%d", key, level)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.Level); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.Level, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	classLevel, err := c.client.GetClassLevel(key, level)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, classLevel)
	return classLevel, nil
}

// GetProficiency returns cached proficiency or fetches from API
func (c *CachedClient) GetProficiency(key string) (*entities.Proficiency, error) {
	cacheKey := fmt.Sprintf("proficiency:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.Proficiency); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.Proficiency, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	proficiency, err := c.client.GetProficiency(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, proficiency)
	return proficiency, nil
}

// ListDamageTypes returns cached damage type list or fetches from API
func (c *CachedClient) ListDamageTypes() ([]*entities.ReferenceItem, error) {
	cacheKey := "list:damage-types"
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	damageTypes, err := c.client.ListDamageTypes()
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, damageTypes)
	return damageTypes, nil
}

// GetDamageType returns cached damage type or fetches from API
func (c *CachedClient) GetDamageType(key string) (*entities.DamageType, error) {
	cacheKey := fmt.Sprintf("damage-type:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.DamageType); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.DamageType, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	damageType, err := c.client.GetDamageType(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, damageType)
	return damageType, nil
}

// GetEquipmentCategory returns cached equipment category or fetches from API
func (c *CachedClient) GetEquipmentCategory(key string) (*entities.EquipmentCategory, error) {
	cacheKey := fmt.Sprintf("equipment-category:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.EquipmentCategory); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.EquipmentCategory, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	category, err := c.client.GetEquipmentCategory(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, category)
	return category, nil
}

// ListBackgrounds returns cached backgrounds list or fetches from API
func (c *CachedClient) ListBackgrounds() ([]*entities.ReferenceItem, error) {
	cacheKey := "backgrounds"
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.([]*entities.ReferenceItem); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected []*entities.ReferenceItem, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	backgrounds, err := c.client.ListBackgrounds()
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, backgrounds)
	return backgrounds, nil
}

// GetBackground returns cached background or fetches from API
func (c *CachedClient) GetBackground(key string) (*entities.Background, error) {
	cacheKey := fmt.Sprintf("background:%s", key)
	
	if cached, ok := c.getFromCache(cacheKey); ok {
		if typedResult, ok := cached.(*entities.Background); ok {
			return typedResult, nil
		}
		log.Printf("Cache type mismatch for key %s, expected *entities.Background, got %T", cacheKey, cached)
		// Fall through to API call
	}
	
	// Cache miss - fetch from API
	background, err := c.client.GetBackground(key)
	if err != nil {
		return nil, err
	}
	
	c.storeInCache(cacheKey, background)
	return background, nil
}