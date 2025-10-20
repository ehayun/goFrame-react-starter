package controller

import (
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"

	"tzlev/internal/model"
	"tzlev/internal/repository"
)

type ClassroomController struct {
	classroomRepo *repository.ClassroomRepository
}

func NewClassroomController() *ClassroomController {
	return &ClassroomController{
		classroomRepo: repository.NewClassroomRepository(),
	}
}

// GetClassrooms retrieves classrooms based on query parameters
func (c *ClassroomController) GetClassrooms(r *ghttp.Request) {
	ctx := gctx.New()

	// Get query parameters
	academicYear := r.Get("academic_year").String()
	schoolIDStr := r.Get("school_id").String()
	teacherIDStr := r.Get("teacher_id").String()

	var classrooms []model.Classroom
	var err error

	// Filter by academic year if provided
	if academicYear != "" {
		classrooms, err = c.classroomRepo.FindByAcademicYear(ctx, academicYear)
	} else if schoolIDStr != "" {
		schoolID, parseErr := strconv.ParseInt(schoolIDStr, 10, 64)
		if parseErr != nil {
			r.Response.WriteJson(g.Map{
				"success": false,
				"message": "Invalid school_id",
			})
			return
		}
		classrooms, err = c.classroomRepo.FindBySchoolID(ctx, schoolID)
	} else if teacherIDStr != "" {
		teacherID, parseErr := strconv.ParseInt(teacherIDStr, 10, 64)
		if parseErr != nil {
			r.Response.WriteJson(g.Map{
				"success": false,
				"message": "Invalid teacher_id",
			})
			return
		}
		classrooms, err = c.classroomRepo.FindByTeacherID(ctx, teacherID)
	} else {
		// Get all classrooms with pagination
		offset := r.Get("offset", 0).Int()
		limit := r.Get("limit", 50).Int()
		classrooms, err = c.classroomRepo.List(ctx, offset, limit)
	}

	if err != nil {
		g.Log().Error(ctx, "Error getting classrooms:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to retrieve classrooms",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":    true,
		"classrooms": classrooms,
	})
}

// GetClassroom retrieves a specific classroom by ID
func (c *ClassroomController) GetClassroom(r *ghttp.Request) {
	ctx := gctx.New()

	idStr := r.Get("id").String()
	if idStr == "" {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Classroom ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Invalid classroom ID",
		})
		return
	}

	classroom, err := c.classroomRepo.FindByID(ctx, id)
	if err != nil {
		g.Log().Error(ctx, "Error getting classroom:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Classroom not found",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":   true,
		"classroom": classroom,
	})
}

// CreateClassroom creates a new classroom
func (c *ClassroomController) CreateClassroom(r *ghttp.Request) {
	ctx := gctx.New()

	var classroom model.Classroom
	if err := r.Parse(&classroom); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}

	if err := c.classroomRepo.Create(ctx, &classroom); err != nil {
		g.Log().Error(ctx, "Error creating classroom:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to create classroom",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":   true,
		"message":   "Classroom created successfully",
		"classroom": classroom,
	})
}

// UpdateClassroom updates an existing classroom
func (c *ClassroomController) UpdateClassroom(r *ghttp.Request) {
	ctx := gctx.New()

	idStr := r.Get("id").String()
	if idStr == "" {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Classroom ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Invalid classroom ID",
		})
		return
	}

	var classroom model.Classroom
	if err := r.Parse(&classroom); err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}

	classroom.ID = id

	if err := c.classroomRepo.Update(ctx, &classroom); err != nil {
		g.Log().Error(ctx, "Error updating classroom:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to update classroom",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success":   true,
		"message":   "Classroom updated successfully",
		"classroom": classroom,
	})
}

// DeleteClassroom deletes a classroom
func (c *ClassroomController) DeleteClassroom(r *ghttp.Request) {
	ctx := gctx.New()

	idStr := r.Get("id").String()
	if idStr == "" {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Classroom ID is required",
		})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Invalid classroom ID",
		})
		return
	}

	if err := c.classroomRepo.Delete(ctx, id); err != nil {
		g.Log().Error(ctx, "Error deleting classroom:", err)
		r.Response.WriteJson(g.Map{
			"success": false,
			"message": "Failed to delete classroom",
		})
		return
	}

	r.Response.WriteJson(g.Map{
		"success": true,
		"message": "Classroom deleted successfully",
	})
}
