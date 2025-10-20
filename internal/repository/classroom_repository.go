package repository

import (
	"context"
	"time"

	"tzlev/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type ClassroomRepository struct{}

func NewClassroomRepository() *ClassroomRepository {
	return &ClassroomRepository{}
}

func (r *ClassroomRepository) Create(ctx context.Context, classroom *model.Classroom) error {
	classroom.InsertedAt = time.Now()
	classroom.UpdatedAt = time.Now()

	_, err := g.DB().Model("classrooms").Ctx(ctx).Insert(classroom)
	return err
}

func (r *ClassroomRepository) FindByID(ctx context.Context, id int64) (*model.Classroom, error) {
	var classroom model.Classroom
	err := g.DB().Model("classrooms").Ctx(ctx).
		Where("id = ?", id).
		Scan(&classroom)

	if err != nil {
		return nil, err
	}
	return &classroom, nil
}

func (r *ClassroomRepository) FindByCode(ctx context.Context, code string) (*model.Classroom, error) {
	var classroom model.Classroom
	err := g.DB().Model("classrooms").Ctx(ctx).
		Where("code = ?", code).
		Scan(&classroom)

	if err != nil {
		return nil, err
	}
	return &classroom, nil
}

func (r *ClassroomRepository) FindByAcademicYear(ctx context.Context, academicYear string) ([]model.Classroom, error) {
	var classrooms []model.Classroom
	err := g.DB().Model("classrooms").Ctx(ctx).
		Where("academic_year = ?", academicYear).
		Order("order_id ASC, classroom_name ASC").
		Scan(&classrooms)

	return classrooms, err
}

func (r *ClassroomRepository) FindBySchoolID(ctx context.Context, schoolID int64) ([]model.Classroom, error) {
	var classrooms []model.Classroom
	err := g.DB().Model("classrooms").Ctx(ctx).
		Where("school_id = ?", schoolID).
		Order("order_id ASC, classroom_name ASC").
		Scan(&classrooms)

	return classrooms, err
}

func (r *ClassroomRepository) FindByTeacherID(ctx context.Context, teacherID int64) ([]model.Classroom, error) {
	var classrooms []model.Classroom
	err := g.DB().Model("classrooms").Ctx(ctx).
		Where("teacher_id = ?", teacherID).
		Order("order_id ASC, classroom_name ASC").
		Scan(&classrooms)

	return classrooms, err
}

func (r *ClassroomRepository) Update(ctx context.Context, classroom *model.Classroom) error {
	classroom.UpdatedAt = time.Now()

	_, err := g.DB().Model("classrooms").Ctx(ctx).
		Where("id = ?", classroom.ID).
		Update(classroom)
	return err
}

func (r *ClassroomRepository) Delete(ctx context.Context, id int64) error {
	_, err := g.DB().Model("classrooms").Ctx(ctx).
		Where("id = ?", id).
		Delete()
	return err
}

func (r *ClassroomRepository) GetDistinctAcademicYears(ctx context.Context) ([]string, error) {
	var results []struct {
		AcademicYear string `json:"academic_year"`
	}
	err := g.DB().Model("classrooms").Ctx(ctx).
		Fields("DISTINCT academic_year").
		Where("academic_year IS NOT NULL AND academic_year != ''").
		Order("academic_year DESC").
		Scan(&results)

	if err != nil {
		return nil, err
	}

	// Extract just the values
	academicYears := make([]string, len(results))
	for i, result := range results {
		academicYears[i] = result.AcademicYear
	}

	return academicYears, nil
}

func (r *ClassroomRepository) List(ctx context.Context, offset, limit int) ([]model.Classroom, error) {
	var classrooms []model.Classroom
	err := g.DB().Model("classrooms").Ctx(ctx).
		Order("order_id ASC, classroom_name ASC").
		Offset(offset).
		Limit(limit).
		Scan(&classrooms)

	return classrooms, err
}
