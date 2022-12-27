package dnd5e

import "github.com/fadedpez/dnd5e-api/entities"

type Interface interface {
	ListRaces() ([]*entities.Race, error)
	GetRace(key string) (*entities.Race, error)
}
