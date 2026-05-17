package service

import (
	"context"
	"fmt"
	"user-service/internal/domain"
)

type UserService struct {
	repo domain.Storage
}

func NewUserService(repo domain.Storage) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, id, name string, age int) error {
	u := domain.User{
		ID:   id,
		Name: name,
		Age:  age,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return fmt.Errorf("CreateUser: %w", err)
	}
	return nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("GetUser: %w", err)
	}
	return user, nil
}

func (s *UserService) ListByAgeRange(ctx context.Context, min, max int) ([]domain.User, error) {
	result, err := s.repo.ListByAgeRange(ctx, min, max)
	if err != nil {
		return nil, fmt.Errorf("ListByAgeRange: %w", err)
	}
	return result, nil
}
