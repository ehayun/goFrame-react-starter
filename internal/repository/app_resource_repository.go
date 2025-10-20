package repository

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"tzlev/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type AppResourceRepository struct{}

func NewAppResourceRepository() *AppResourceRepository {
	return &AppResourceRepository{}
}

// generateRandomID generates a random 10-character alphanumeric string
func generateRandomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rng.Intn(len(charset))]
	}
	return string(result)
}

func (r *AppResourceRepository) Create(ctx context.Context, resource *model.AppResource) error {
	if resource.Id == "" {
		resource.Id = generateRandomID()
	}

	// Check if a resource with the same name already exists
	existingResource, err := r.FindByName(ctx, resource.Name)
	if err == nil && existingResource != nil {
		return fmt.Errorf("a resource with name '%s' already exists", resource.Name)
	}

	_, err = g.DB().Model("app_resources").Ctx(ctx).Insert(resource)
	return err
}

func (r *AppResourceRepository) FindByID(ctx context.Context, id string) (*model.AppResource, error) {
	var resource model.AppResource
	err := g.DB().Model("app_resources").Ctx(ctx).
		Where("id = ?", id).
		Scan(&resource)

	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *AppResourceRepository) FindByName(ctx context.Context, name string) (*model.AppResource, error) {
	var resource model.AppResource
	err := g.DB().Model("app_resources").Ctx(ctx).
		Where("name = ?", name).
		Scan(&resource)

	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *AppResourceRepository) Update(ctx context.Context, resource *model.AppResource) error {
	_, err := g.DB().Model("app_resources").Ctx(ctx).
		Where("id = ?", resource.Id).
		Update(resource)
	return err
}

func (r *AppResourceRepository) Delete(ctx context.Context, id string) error {
	_, err := g.DB().Model("app_resources").Ctx(ctx).
		Where("id = ?", id).
		Delete()
	return err
}

func (r *AppResourceRepository) List(ctx context.Context, offset, limit int) ([]model.AppResource, error) {
	var resources []model.AppResource
	err := g.DB().Model("app_resources").Ctx(ctx).
		Order("name ASC").
		Offset(offset).
		Limit(limit).
		Scan(&resources)

	return resources, err
}

func (r *AppResourceRepository) ListAll(ctx context.Context) ([]model.AppResource, error) {
	var resources []model.AppResource
	err := g.DB().Model("app_resources").Ctx(ctx).
		Order("name ASC").
		Scan(&resources)

	return resources, err
}
