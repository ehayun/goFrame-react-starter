package service

import (
	"context"
	"fmt"
	"time"

	"tzlev/internal/cache"
	"tzlev/internal/model"
	"tzlev/internal/repository"
)

type UserService struct {
	userRepo     *repository.UserRepository
	cacheManager *cache.CacheManager
}

func NewUserService() *UserService {
	return &UserService{
		userRepo:     repository.NewUserRepository(),
		cacheManager: cache.NewCacheManager(),
	}
}

func (s *UserService) GetUserByZehut(ctx context.Context, zehut string) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%s", zehut)

	// Try to get from cache first
	var user model.User
	err := s.cacheManager.Get(ctx, cacheKey, &user)
	if err == nil {
		return &user, nil
	}

	// If not in cache, get from database
	dbUser, err := s.userRepo.FindByZehut(ctx, zehut)
	if err != nil {
		return nil, err
	}

	// Store in cache for 5 minutes
	_ = s.cacheManager.Set(ctx, cacheKey, dbUser, 5*time.Minute)

	return dbUser, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

func (s *UserService) UpdateUser(ctx context.Context, user *model.User) error {
	// Update in database
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", user.Zehut)
	_ = s.cacheManager.Delete(ctx, cacheKey)

	return nil
}
