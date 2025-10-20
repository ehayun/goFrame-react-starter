package controller

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"tzlev/internal/redis"
	"tzlev/internal/repository"
)

type AcademicYearController struct {
	classroomRepo *repository.ClassroomRepository
}

func NewAcademicYearController() *AcademicYearController {
	return &AcademicYearController{
		classroomRepo: repository.NewClassroomRepository(),
	}
}

// GetAcademicYear retrieves the user's saved academic year from Redis
func (c *AcademicYearController) GetAcademicYear(r *ghttp.Request) {
	ctx := gctx.New()

	// Get user zehut from session
	userZehut, _ := r.Session.Get("user_zehut")
	if userZehut.String() == "" {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Get academic year from Redis
	key := "tzlev:user:" + userZehut.String() + ":academic_year"
	academicYear, err := redis.Client.Get(ctx, key).Result()
	if err != nil {
		// If key doesn't exist, return empty
		if err.Error() == "redis: nil" {
			r.Response.WriteJson(g.Map{
				"success":      true,
				"academicYear": "",
			})
			return
		}

		g.Log().Error(ctx, "Error getting academic year from Redis:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to retrieve academic year",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":      true,
		"academicYear": academicYear,
	})
}

// SetAcademicYear saves the user's academic year to Redis
func (c *AcademicYearController) SetAcademicYear(r *ghttp.Request) {
	ctx := gctx.New()

	// Get user zehut from session
	userZehut, _ := r.Session.Get("user_zehut")
	if userZehut.String() == "" {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "User not authenticated",
		})
		return
	}

	// Parse request body
	var request struct {
		AcademicYear string `json:"academicYear"`
	}

	if err := r.Parse(&request); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}

	// Validate academic year
	if request.AcademicYear == "" {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Academic year is required",
		})
		return
	}

	// Save academic year to Redis with 1 year expiration
	key := "tzlev:user:" + userZehut.String() + ":academic_year"
	err := redis.Client.Set(ctx, key, request.AcademicYear, 365*24*time.Hour).Err()
	if err != nil {
		g.Log().Error(ctx, "Error saving academic year to Redis:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to save academic year",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Academic year saved successfully",
	})
}

// GetAcademicYearsList retrieves all available academic years from classrooms table
func (c *AcademicYearController) GetAcademicYearsList(r *ghttp.Request) {
	ctx := gctx.New()

	academicYears, err := c.classroomRepo.GetDistinctAcademicYears(ctx)
	if err != nil {
		g.Log().Error(ctx, "Error getting academic years list:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to retrieve academic years",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":       true,
		"academicYears": academicYears,
	})
}

// GetCurrentAcademicYear retrieves the user's current academic year from Redis
// This is a general function that can be used by other parts of the application
func GetCurrentAcademicYear(userZehut string) (string, error) {
	ctx := gctx.New()

	if userZehut == "" {
		return "", fmt.Errorf("user zehut is required")
	}

	// Get academic year from Redis
	key := "tzlev:user:" + userZehut + ":academic_year"
	academicYear, err := redis.Client.Get(ctx, key).Result()
	if err != nil {
		// If key doesn't exist, return empty string
		if err.Error() == "redis: nil" {
			return "", nil
		}
		return "", err
	}

	return academicYear, nil
}

// GetCurrentAcademicYearFromRequest retrieves the current academic year from the request session
// This is a convenience function for controllers that have access to the request
func (c *AcademicYearController) GetCurrentAcademicYearFromRequest(r *ghttp.Request) (string, error) {
	// Get user zehut from session
	userZehut, _ := r.Session.Get("user_zehut")
	if userZehut.String() == "" {
		return "", fmt.Errorf("user not authenticated")
	}

	return GetCurrentAcademicYear(userZehut.String())
}
