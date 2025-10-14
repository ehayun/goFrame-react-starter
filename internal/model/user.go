package model

import (
	"time"
)

// User matches the existing users table in the database
type User struct {
	Zehut              string     `json:"zehut" orm:"zehut"`                             // Primary key
	LastName           string     `json:"last_name" orm:"last_name"`
	FirstName          string     `json:"first_name" orm:"first_name"`
	Mobile             string     `json:"mobile,omitempty" orm:"mobile"`
	Avatar             string     `json:"avatar,omitempty" orm:"avatar"`
	Role               string     `json:"role,omitempty" orm:"role"`
	IsAdmin            bool       `json:"is_admin" orm:"is_admin"`
	Email              string     `json:"email,omitempty" orm:"email"`
	DateOfBirth        *time.Time `json:"date_of_birth,omitempty" orm:"date_of_birth"`
	HashedPassword     string     `json:"-" orm:"hashed_password"` // Don't serialize password
	ConfirmedAt        *time.Time `json:"confirmed_at,omitempty" orm:"confirmed_at"`
	InsertedAt         time.Time  `json:"inserted_at" orm:"inserted_at"`
	UpdatedAt          time.Time  `json:"updated_at" orm:"updated_at"`
	MeravId            string     `json:"merav_id,omitempty" orm:"merav_id"`
	TagNo              string     `json:"tag_no,omitempty" orm:"tag_no"`
	Phone              string     `json:"phone,omitempty" orm:"phone"`
	RoleDescription    string     `json:"role_description,omitempty" orm:"role_description"`
	Remarks            string     `json:"remarks,omitempty" orm:"remarks"`
	DailyAllowance     bool       `json:"daily_allowance" orm:"daily_allowance"`
	MonthlyAllowance   bool       `json:"monthly_allowance" orm:"monthly_allowance"`
	TravelConsiderate  bool       `json:"travel_considerate" orm:"travel_considerate"`
	MeravMifal         string     `json:"merav_mifal,omitempty" orm:"merav_mifal"`
	StartOfWork        *time.Time `json:"start_of_work,omitempty" orm:"start_of_work"`
	ManagerId          string     `json:"manager_id,omitempty" orm:"manager_id"`
	Address            string     `json:"address,omitempty" orm:"address"`
	City               string     `json:"city,omitempty" orm:"city"`
	Zipcode            string     `json:"zipcode,omitempty" orm:"zipcode"`
	Allowance          int        `json:"allowance" orm:"allowance"`
	PaymentPerHour     int        `json:"payment_per_hour" orm:"payment_per_hour"`
	HoursCount         bool       `json:"hours_count" orm:"hours_count"`
	MonthlyTicket      int        `json:"monthly_ticket" orm:"monthly_ticket"`
	UnitCode           string     `json:"unit_code,omitempty" orm:"unit_code"`
	Id                 *int64     `json:"id,omitempty" orm:"id"`
	IsFreezed          *bool      `json:"is_freezed,omitempty" orm:"is_freezed"`
	FullName           string     `json:"full_name,omitempty" orm:"full_name"`
	CanEncrypt         *bool      `json:"can_encrypt,omitempty" orm:"can_encrypt"`
	EncryptPassword    string     `json:"-" orm:"encrypt_password"` // Don't serialize
	CanUpdateDocs      bool       `json:"can_update_docs" orm:"can_update_docs"`
	EmailToken         string     `json:"-" orm:"email_token"` // Don't serialize token
	Color              string     `json:"color,omitempty" orm:"color"`
	LicenceNumber      string     `json:"licence_number,omitempty" orm:"licence_number"`
	CanUpdateDisciplines bool     `json:"can_update_disciplines" orm:"can_update_disciplines"`
	ConfirmAttendance  string     `json:"confirm_attendance,omitempty" orm:"confirm_attendance"`
	AllowedEditUsers   bool       `json:"allowed_edit_users" orm:"allowed_edit_users"`
	IsTechnical        bool       `json:"is_technical" orm:"is_technical"`
	CanSeeSensitiveDocs bool      `json:"can_see_sensitive_docs" orm:"can_see_sensitive_docs"`
}
