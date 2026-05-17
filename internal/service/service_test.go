package service_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"testing"
	"user-service/internal/domain"
	"user-service/internal/service"
	"user-service/internal/storage"
)

var (
	user_1 = domain.User{
		ID:   "user-1",
		Name: "Alice",
		Age:  18,
	}
	user_2 = domain.User{
		ID:   "user-2",
		Name: "Bob",
		Age:  24,
	}
	user_3 = domain.User{
		ID:   "user-3",
		Name: "Chuck",
		Age:  68,
	}
)

func TestCreateUser(t *testing.T) {
	repo := storage.NewInMemoryStorage()
	svc := service.NewUserService(repo)
	ctx := context.Background()

	tests := []struct {
		testName string
		id       string
		name     string
		age      int
		wantErr  error
	}{
		{
			testName: "Success creation",
			id:       user_1.ID,
			name:     user_1.Name,
			age:      user_1.Age,
			wantErr:  nil,
		},
		{
			testName: "Success creation",
			id:       user_2.ID,
			name:     user_2.Name,
			age:      user_2.Age,
			wantErr:  nil,
		},
		{
			testName: "Error: duplicate id",
			id:       user_1.ID,
			name:     user_1.Name,
			age:      user_1.Age,
			wantErr:  domain.ErrDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := svc.CreateUser(ctx, tt.id, tt.name, tt.age)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateUser() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	repo := storage.NewInMemoryStorage()
	svc := service.NewUserService(repo)
	ctx := context.Background()

	tests := []struct {
		testName string
		id       string
		wantUsr  *domain.User
		wantErr  error
	}{
		{
			testName: "Success adding",
			id:       user_1.ID,
			wantUsr:  &user_1,
			wantErr:  nil,
		},
		{
			testName: "Error: user not found",
			id:       user_2.ID,
			wantUsr:  nil,
			wantErr:  domain.ErrNotFound,
		},
	}

	svc.CreateUser(ctx, user_1.ID, user_1.Name, user_1.Age)

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			user, err := svc.GetUser(ctx, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetUser() err = %v, wantErr = %v", err, tt.wantErr)
			}
			if user != nil && *user != *tt.wantUsr {
				t.Errorf("GetUser() user = %v, wantUsr = %v", user, tt.wantUsr)
			}
		})
	}
}

func TestListByAgeRange(t *testing.T) {
	repo := storage.NewInMemoryStorage()
	svc := service.NewUserService(repo)
	ctx := context.Background()

	svc.CreateUser(ctx, user_1.ID, user_1.Name, user_1.Age)
	svc.CreateUser(ctx, user_2.ID, user_2.Name, user_2.Age)
	svc.CreateUser(ctx, user_3.ID, user_3.Name, user_3.Age)

	tests := []struct {
		testName string
		min      int
		max      int
		wantRes  []domain.User
		wantErr  error
	}{
		{
			testName: "Success result",
			min:      1,
			max:      70,
			wantRes:  []domain.User{user_1, user_2, user_3},
			wantErr:  nil,
		},
		{
			testName: "Error: user not found",
			min:      0,
			max:      0,
			wantRes:  nil,
			wantErr:  domain.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			res, err := svc.ListByAgeRange(ctx, tt.min, tt.max)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ListByAgeRange() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !slices.Equal(res, tt.wantRes) {
				t.Errorf("ListByAgeRange() result = %v, wantRes = %v", res, tt.wantRes)
			}
		})
	}
}

func TestInMemoryStorage_Concurrency(t *testing.T) {
	repo := storage.NewInMemoryStorage()
	svc := service.NewUserService(repo)
	ctx := context.Background()

	const workers = 100
	const opsPerWorker = 100

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := range workers {
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < opsPerWorker; j++ {
				id := fmt.Sprintf("u-%d-%d", workerID, j)
				// user := domain.User{ID: id, Name: "Test"}

				_ = svc.CreateUser(ctx, id, "Test", 20)
				_, _ = svc.GetUser(ctx, id)
			}
		}(i)
	}
	wg.Wait()

}
