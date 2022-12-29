package dnd5e

import (
	"net/http"

	"github.com/fadedpez/dnd5e-api/entities"
)

type Interface interface {
	ListRaces() ([]*entities.Race, error)
	GetRace(key string) (*entities.Race, error)
}

type httpIface interface {
	Get(url string) (*http.Response, error)
}
