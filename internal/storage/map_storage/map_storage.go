package mapstorage

import (
	"auth_service/internal/domain/models"
	"auth_service/internal/storage"
	"context"
	"sync"
)

type MapStorage struct {
	m       sync.RWMutex
	users   map[string]models.User
	apps 	map[int64]models.App
	counter int64
}

func New() *MapStorage {
	apps := make(map[int64]models.App)
	apps[1] = models.App{
		ID:1,
		Secret:"sbebraboy",
	}
	return &MapStorage{
		users: make(map[string]models.User),
		apps: apps,
	}
}

func (m *MapStorage) User(ctx context.Context, email string) (models.User, error) {
	user, ok := m.isUserExists(email)
	if !ok {
		return user, storage.ErrUserNotFound
	}

	return user, nil
}

func (m *MapStorage) isUserExists(email string) (models.User, bool) {
	m.m.RLock()
	u, ok := m.users[email]
	m.m.RUnlock()
	return u, ok
}

func (m *MapStorage) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	_, ok := m.isUserExists(email)
	if ok {
		return 0, storage.ErrUserExists
	}

	m.m.Lock()
	defer m.m.Unlock()

	m.counter++
	m.users[email] = models.User{
		ID:       m.counter,
		Email:    email,
		PassHash: passHash,
	}

	return m.counter, nil
}

func (m *MapStorage) App(ctx context.Context, appID int) (models.App, error) {
	v, ok := m.apps[int64(appID)]
	if !ok {
		return models.App{}, storage.ErrAppNotFound
	}
	return v, nil
}
