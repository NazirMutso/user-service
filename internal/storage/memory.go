package storage

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"user-service/internal/domain"
)

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

func (s *InMemoryStorage) GetByID(ctx context.Context, id string) (*domain.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return &val, nil
}

func (s *InMemoryStorage) ListByAgeRange(ctx context.Context, min, max int) ([]domain.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]domain.User, 0, len(s.data))
	for _, v := range s.data {
		if v.Age >= min && v.Age <= max {
			result = append(result, v)
		}
	}

	if len(result) == 0 {
		return nil, domain.ErrNotFound
	}

	slices.SortFunc(result, func(a, b domain.User) int {
		return cmp.Compare(a.Age, b.Age)
	})

	return result, nil
}
