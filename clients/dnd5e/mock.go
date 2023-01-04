package dnd5e

import (
	"net/http"

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
