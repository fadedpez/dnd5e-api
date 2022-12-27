package dnd5e

import (
	"net/http"

	"github.com/fadedpez/dnd5e-api/entities"

	"github.com/stretchr/testify/mock"
)

type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*http.Response), args.Error(1)
}

type mockDND5eAPI struct {
	mock.Mock
}

func (m *mockDND5eAPI) ListRaces() ([]*entities.Race, error) {
	args := m.Called()
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*entities.Race), args.Error(1)
}

func (m *mockDND5eAPI) GetRace(key string) (*entities.Race, error) {
	args := m.Called(key)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entities.Race), args.Error(1)
}
