package storage

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"user-service/internal/domain"
)

// Compile-time interface check
var _ domain.Storage = (*InMemoryStorage)(nil)

type InMemoryStorage struct {
	mu   sync.RWMutex
	data map[string]domain.User
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{data: make(map[string]domain.User)}
}

func (s *InMemoryStorage) Create(ctx context.Context, u domain.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[u.ID]; ok {
		return domain.ErrDuplicate
	}
	s.data[u.ID] = u

	return nil
}

func (s *InMemoryStorage) GetByID(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	select {
	case <-ctx.Done():
		return user, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[id]
	if !ok {
		return user, domain.ErrNotFound
	}
	return val, nil
}

func (s *InMemoryStorage) ListByAgeRange(ctx context.Context, min, max int) (usr []domain.User, err error) {
	select {
	case <-ctx.Done():
		return usr, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	usr = make([]domain.User, 0, len(s.data))
	for _, v := range s.data {
		if v.Age >= min && v.Age <= max {
			usr = append(usr, v)
		}
	}

	if len(usr) == 0 {
		return usr, nil
	}

	slices.SortFunc(usr, func(a, b domain.User) int {
		return cmp.Compare(a.Age, b.Age)
	})

	return usr, nil
}
