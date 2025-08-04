package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null column:password" json:"-"`
}
type UserCreate struct {
	Username string `json:"username" example:"chayan_das"`
	Password string `json:"password" example:"password123"`
}
type UploadDocument struct {
	ID             uint           `gorm:"primaryKey" json:"-"`
	StudentID      uint           `gorm:"not null" json:"-"`
	Name           string         `json:"name"`
	URL            string         `json:"url"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	StudentProfile StudentProfile `gorm:"foreignKey:StudentID" json:"-"`
}

const (
	AddressTypePermanent = "permanent"
	AddressTypeCurrent   = "current"
)

type Address struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	StudentID uint      `gorm:"not null" json:"-"`
	Type      string    `gorm:"type:varchar(20)" json:"type"` // permanent, current
	Street    string    `json:"street"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Pincode   string    `gorm:"type:varchar(10)" json:"pincode"`
	Country   string    `json:"country"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type StudentAcademicQualification struct {
	ID            uint      `gorm:"primaryKey" json:"-"`
	StudentID     uint      `gorm:"not null" json:"-"`
	Degree        string    `json:"degree"`
	University    string    `json:"university"`
	YearOfPassing int       `json:"year_of_passing"`
	Grade         string    `json:"grade"`
	Course        string    `json:"course"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type StudentProfile struct {
	ID               uint                           `gorm:"primaryKey" json:"-"`
	UserID           uint                           `gorm:"not null;" json:"-"`
	FullName         string                         `gorm:"not null" json:"full_name"`
	DateOfBirth      string                         `json:"date_of_birth" gorm:"type:varchar(10)"` // Format: DD-MM-YYYY
	Gender           string                         `gorm:"type:varchar(10)" json:"gender"`
	PhoneNumber      string                         `gorm:"type:varchar(15)" json:"phone_number"`
	Qualification    string                         `gorm:"type:varchar(50)" json:"qualification"`
	Email            string                         `gorm:"type:varchar(100)" json:"email"`
	AadhaarNumber    string                         `gorm:"type:varchar(12)" json:"aadhaar_number"`
	Nationality      string                         `json:"nationality"` // Added Nationality
	Category         string                         `gorm:"type:varchar(20)" json:"category"`
	Income           float64                        `json:"income"`
	IsInternational  bool                           `json:"is_international"` // Flag to mark international students
	CreatedAt        time.Time                      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time                      `gorm:"autoUpdateTime" json:"updated_at"`
	Documents        []UploadDocument               `gorm:"foreignKey:StudentID" json:"documents"`
	EducationHistory []StudentAcademicQualification `gorm:"foreignKey:StudentID" json:"education_history"` // Academic qualifications
	Addresses        []Address                      `gorm:"foreignKey:StudentID" json:"addresses"`         // List of addresses
}

// Application represents a scholarship application
type Application struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"not null" json:"user_id"`
	SchemeID         uint           `gorm:"not null" json:"scheme_id"`
	StudentProfileID uint           `gorm:"not null" json:"student_profile_id"`
	IsDraft          bool           `gorm:"default:true" json:"is_draft"`
	Verified         bool           `gorm:"default:false" json:"verified"`
	SubmittedAt      *time.Time     `json:"submitted_at,omitempty"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	User             User           `gorm:"foreignKey:UserID" json:"user"`
	Scheme           Scheme         `gorm:"foreignKey:SchemeID" json:"-"`
	StudentProfile   StudentProfile `gorm:"foreignKey:StudentProfileID" json:"student_profile"`
	Status           string         `gorm:"type:varchar(20)" json:"status"` // status of the application (submitted, under review, approved, rejected)
}

type DocumentInput struct {
	Name string `json:"name" example:"Aadhaar Card"`
	URL  string `json:"url" example:"https://example.com/documents/aadhaar_card.pdf"`
}

type AddressInput struct {
	Type    string `json:"type" example:"permanent"` // e.g., permanent, current
	Street  string `json:"street" example:"123 Main St"`
	City    string `json:"city" example:"Kolkata"`
	State   string `json:"state" example:"West Bengal"`
	Pincode string `json:"pincode" example:"700001"`
	Country string `json:"country" example:"India"`
}

type EducationHistoryInput struct {
	Degree        string `json:"degree" example:"Bachelor of Science"`
	University    string `json:"university" example:"University of Calcutta"`
	YearOfPassing int    `json:"year_of_passing" example:"2022"`
	Grade         string `json:"grade" example:"First Class"`
	Course        string `json:"course" example:"Physics"`
}

type StudentProfileInput struct {
	FullName         string                  `json:"full_name" example:"Chayan Das"`
	DateOfBirth      string                  `json:"date_of_birth" example:"01/10/2004"` // Format: DD/MM/YYYY
	Gender           string                  `json:"gender" example:"Male" `
	PhoneNumber      string                  `json:"phone_number" example:"+919876543210"`
	Qualification    string                  `json:"qualification" example:"Bachelor of Science"`
	Email            string                  `json:"email" example:"das@gmail.com"`
	AadhaarNumber    string                  `json:"aadhaar_number" example:"123456789012"`
	Nationality      string                  `json:"nationality" example:"Indian"` // Added Nationality
	Category         string                  `json:"category" example:"General"`
	Income           float64                 `json:"income" example:"250000.00"`
	IsInternational  bool                    `json:"is_international" example:"false"` // Flag to mark international students
	Documents        []DocumentInput         `json:"documents"`
	Addresses        []AddressInput          `json:"addresses"`
	EducationHistory []EducationHistoryInput `json:"education_history"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	// Validate address type
	if a.Type != AddressTypePermanent && a.Type != AddressTypeCurrent {
		return fmt.Errorf("invalid address type: %s", a.Type)
	}

	var count int64
	if err := tx.Model(&Address{}).
		Where("student_id = ? AND type = ?", a.StudentID, a.Type).
		Count(&count).Error; err != nil {
		return fmt.Errorf("error counting existing addresses: %w", err)
	}

	fmt.Printf("Existing count for type '%s': %d\n", a.Type, count)

	if count >= 1 {
		return fmt.Errorf("a student can have only one '%s' address", a.Type)
	}
	return nil
}

func (a *Application) BeforeCreate(tx *gorm.DB) (err error) {
	if a.IsDraft {
		a.SubmittedAt = nil
	} else {
		if a.SubmittedAt == nil {
			now := time.Now()
			a.SubmittedAt = &now
		}
	}
	return
}

type InitApplicationRequest struct {
	SchemeID uint `json:"scheme_id" binding:"required"`
}

type SubmitExistingApplicationRequest struct {
	ApplicationID uint `json:"application_id" binding:"required"`
}

type Validator interface {
	Name() string
	Validate(app Application) error
}
type ValidationChain interface {
	AddValidator(validator Validator) ValidationChain
	Validate(app Application) error
}

type validationChain struct {
	validators []Validator
}
