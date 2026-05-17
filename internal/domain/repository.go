package domain

import "context"

type Storage interface {
	Create(ctx context.Context, u User) error
	GetByID(ctx context.Context, id string) (User, error)
	ListByAgeRange(ctx context.Context, min, max int) ([]User, error)
}
