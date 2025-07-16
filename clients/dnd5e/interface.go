package dnd5e

import (
	"net/http"

	"github.com/fadedpez/dnd5e-api/entities"
)

type Interface interface {
	ListRaces() ([]*entities.ReferenceItem, error)
	GetRace(key string) (*entities.Race, error)
	ListEquipment() ([]*entities.ReferenceItem, error)
	GetEquipment(key string) (EquipmentInterface, error)
	ListClasses() ([]*entities.ReferenceItem, error)
	GetClass(key string) (*entities.Class, error)
	ListSpells(input *ListSpellsInput) ([]*entities.ReferenceItem, error)
	GetSpell(key string) (*entities.Spell, error)
	ListFeatures() ([]*entities.ReferenceItem, error)
	GetFeature(key string) (*entities.Feature, error)
	ListSkills() ([]*entities.ReferenceItem, error)
	GetSkill(key string) (*entities.Skill, error)
	ListMonsters() ([]*entities.ReferenceItem, error)
	ListMonstersWithFilter(input *ListMonstersInput) ([]*entities.ReferenceItem, error)
	GetMonster(key string) (*entities.Monster, error)
	GetClassLevel(key string, level int) (*entities.Level, error)
	GetProficiency(key string) (*entities.Proficiency, error)
	ListDamageTypes() ([]*entities.ReferenceItem, error)
	GetDamageType(key string) (*entities.DamageType, error)
	GetEquipmentCategory(key string) (*entities.EquipmentCategory, error)
	ListBackgrounds() ([]*entities.ReferenceItem, error)
	GetBackground(key string) (*entities.Background, error)
}

type httpIface interface {
	Get(url string) (*http.Response, error)
}

type EquipmentInterface interface {
	GetType() string
}
