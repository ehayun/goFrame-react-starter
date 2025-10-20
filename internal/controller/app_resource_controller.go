package controller

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"tzlev/internal/model"
	"tzlev/internal/repository"
)

type AppResourceController struct {
	resourceRepo *repository.AppResourceRepository
}

func NewAppResourceController() *AppResourceController {
	return &AppResourceController{
		resourceRepo: repository.NewAppResourceRepository(),
	}
}

// GetAppResources retrieves all app resources
func (c *AppResourceController) GetAppResources(r *ghttp.Request) {
	ctx := gctx.New()

	resources, err := c.resourceRepo.ListAll(ctx)
	if err != nil {
		g.Log().Error(ctx, "Error getting app resources:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to retrieve app resources",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":   true,
		"resources": resources,
	})
}

// GetAppResource retrieves a specific app resource by ID
func (c *AppResourceController) GetAppResource(r *ghttp.Request) {
	ctx := gctx.New()

	id := r.Get("id").String()
	if id == "" {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Resource ID is required",
		})
		return
	}

	resource, err := c.resourceRepo.FindByID(ctx, id)
	if err != nil {
		g.Log().Error(ctx, "Error getting app resource:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Resource not found",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":  true,
		"resource": resource,
	})
}

// CreateAppResource creates a new app resource
func (c *AppResourceController) CreateAppResource(r *ghttp.Request) {
	ctx := gctx.New()

	var resource model.AppResource
	if err := r.Parse(&resource); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}

	if err := c.resourceRepo.Create(ctx, &resource); err != nil {
		g.Log().Error(ctx, "Error creating app resource:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":  true,
		"message":  "App resource created successfully",
		"resource": resource,
	})
}

// UpdateAppResource updates an existing app resource
func (c *AppResourceController) UpdateAppResource(r *ghttp.Request) {
	ctx := gctx.New()

	id := r.Get("id").String()
	if id == "" {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Resource ID is required",
		})
		return
	}

	var resource model.AppResource
	if err := r.Parse(&resource); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}

	resource.Id = id

	if err := c.resourceRepo.Update(ctx, &resource); err != nil {
		g.Log().Error(ctx, "Error updating app resource:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to update app resource",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":  true,
		"message":  "App resource updated successfully",
		"resource": resource,
	})
}

// DeleteAppResource deletes an app resource
func (c *AppResourceController) DeleteAppResource(r *ghttp.Request) {
	ctx := gctx.New()

	id := r.Get("id").String()
	if id == "" {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Resource ID is required",
		})
		return
	}

	if err := c.resourceRepo.Delete(ctx, id); err != nil {
		g.Log().Error(ctx, "Error deleting app resource:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to delete app resource",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "App resource deleted successfully",
	})
}
