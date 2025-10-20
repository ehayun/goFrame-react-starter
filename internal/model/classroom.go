package model

import (
	"time"
)

// ClassType represents the classroom type enum
type ClassType string

const (
	ClassTypeClassroom ClassType = "classroom"
	ClassTypeGroup     ClassType = "group"
)

// SymbolType represents the symbol type enum
type SymbolType string

const (
	SymbolTypeClassroom SymbolType = "classroom"
)

// Classroom matches the existing classrooms table in the database
type Classroom struct {
	ID             int64      `json:"id" orm:"id"`
	Code           string     `json:"code,omitempty" orm:"code"`
	ClassroomName  string     `json:"classroom_name,omitempty" orm:"classroom_name"`
	SchoolID       int64      `json:"school_id" orm:"school_id"`
	TeacherID      int64      `json:"teacher_id" orm:"teacher_id"`
	InsertedAt     time.Time  `json:"inserted_at" orm:"inserted_at"`
	UpdatedAt      time.Time  `json:"updated_at" orm:"updated_at"`
	AcademicYear   string     `json:"academic_year,omitempty" orm:"academic_year"`
	ClassroomType  ClassType  `json:"classroom_type" orm:"classroom_type"`
	ClassroomSemel string     `json:"classroom_semel,omitempty" orm:"classroom_semel"`
	SymbolType     SymbolType `json:"symbol_type,omitempty" orm:"symbol_type"`
	OrderID        int        `json:"order_id,omitempty" orm:"order_id"`
	StartFrom      []string   `json:"start_from,omitempty" orm:"start_from"`
	EndTo          []string   `json:"end_to,omitempty" orm:"end_to"`
	ManualStart    []bool     `json:"manual_start,omitempty" orm:"manual_start"`
}
