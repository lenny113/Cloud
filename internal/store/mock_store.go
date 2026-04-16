package store

import (
	model "assignment-2/internal/models"
	"context"
	"errors"
)

type MockStore struct {
	data    map[string]model.Registration
	apiKeys map[string]bool
}

func NewMockStore() *MockStore {
	return &MockStore{
		data:    make(map[string]model.Registration),
		apiKeys: map[string]bool{"ec654fac9599f62e79e2706abef23dfb7c07c08185aa86db4d8695f0b718d1b3": true},
	}
}

func (m *MockStore) APIKeyExists(ctx context.Context, keyHash string) bool {
	return m.apiKeys[keyHash]
}

func (m *MockStore) ApiKeyExists(ctx context.Context, keyHash string) bool {
	return m.apiKeys[keyHash]
}

func (m *MockStore) CountApiPerUser(ctx context.Context, email string) (int, error) {
	return 0, nil
}

func (m *MockStore) CreateApiStorage(ctx context.Context, reg model.Authentication) error {
	return nil
}

func (m *MockStore) DeleteAPIkey(ctx context.Context, apiKey string) error {
	return nil
}

func (m *MockStore) CreateRegistration(ctx context.Context, apiKey string, reg model.Registration) (string, error) {
	id := "test-id"
	reg.ID = id
	m.data[id] = reg
	return id, nil
}

func (m *MockStore) GetRegistration(ctx context.Context, apiKey string, id string) (*model.Registration, error) {
	reg, ok := m.data[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &reg, nil
}

func (m *MockStore) GetAllRegistrations(ctx context.Context, apiKey string) ([]model.Registration, error) {
	var regs []model.Registration
	for _, reg := range m.data {
		regs = append(regs, reg)
	}
	return regs, nil
}

func (m *MockStore) UpdateRegistration(ctx context.Context, apiKey string, id string, reg model.Registration) error {
	if _, ok := m.data[id]; !ok {
		return errors.New("not found")
	}
	reg.ID = id
	m.data[id] = reg
	return nil
}

func (m *MockStore) DeleteRegistration(ctx context.Context, apiKey string, id string) error {
	if _, ok := m.data[id]; !ok {
		return errors.New("not found")
	}
	delete(m.data, id)
	return nil
}

func (m *MockStore) TweakRegistration(ctx context.Context, apiKey string, id string, patch model.RegistrationPatch) error {
	reg, ok := m.data[id]
	if !ok {
		return errors.New("not found")
	}
	if patch.Country != nil {
		reg.Country = *patch.Country
	}
	m.data[id] = reg
	return nil
}
