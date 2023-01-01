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
	ListSpells() ([]*entities.ReferenceItem, error)
	GetSpell(key string) (*entities.Spell, error)
}

type httpIface interface {
	Get(url string) (*http.Response, error)
}

type EquipmentInterface interface {
	GetType() string
}
