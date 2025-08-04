package utils

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/ChayanDass/beneficiary-manager/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Validate iterates through all the validators in the validation chain
// and applies them to the provided application model. If any validator
// fails, it returns an error indicating which validator failed and the
// specific error encountered.
//
// Parameters:
//   - app: A pointer to the Application model that needs to be validated.
//
// Returns:
//   - error: An error if validation fails at any step, wrapping the name
//     of the failing validator and the specific error. Returns nil if all
//     validations pass.
func (vc *validationChain) Validate(app *models.Application) error {
	for _, validator := range vc.validators {
		if err := validator.Validate(app); err != nil {
			return fmt.Errorf("validation failed at %s: %w", validator.Name(), err)
		}
	}
	return nil
}

type BasicProfileValidator struct{}

func (v BasicProfileValidator) Name() string {
	return "BasicProfileValidator"
}

func (v BasicProfileValidator) Validate(app *models.Application) error {
	return validateBasicProfile(app.StudentProfile)
}

type DocumentsValidator struct{}

func (v DocumentsValidator) Name() string {
	return "DocumentsValidator"
}

func (v DocumentsValidator) Validate(app *models.Application) error {
	return validateDocuments(app.StudentProfile.Documents)
}

type EducationValidator struct{}

func (v EducationValidator) Name() string {
	return "EducationValidator"
}

func (v EducationValidator) Validate(app *models.Application) error {
	return validateEducation(app.StudentProfile.EducationHistory)
}

type AddressesValidator struct{}

func (v AddressesValidator) Name() string {
	return "AddressesValidator"
}

func (v AddressesValidator) Validate(app *models.Application) error {
	return validateAddresses(app.StudentProfile.Addresses)
}

// CheckApplicationCompleteness validates the completeness of a student's application.
func CheckApplicationCompleteness(app *models.Application) error {
	return NewValidationChain().
		AddValidator(BasicProfileValidator{}).
		AddValidator(DocumentsValidator{}).
		AddValidator(EducationValidator{}).
		AddValidator(AddressesValidator{}).
		Validate(app)
}

// ===================== validated application =====================
func validateBasicProfile(p models.StudentProfile) error {
	if p.FullName == "" {
		return errors.New("full name is missing")
	}
	if p.Email == "" {
		return errors.New("email is missing")
	}
	if p.AadhaarNumber == "" {
		return errors.New("aadhaar number is missing")
	}
	if p.PhoneNumber == "" {
		return errors.New("phone number is missing")
	}
	if err := validateDOB(p.DateOfBirth); err != nil {
		return err
	}
	if p.Qualification == "" {
		return errors.New("qualification is missing")
	}
	if p.Nationality == "" {
		return errors.New("nationality is missing")
	}
	if p.Category == "" {
		return errors.New("category is missing")
	}
	if p.Income <= 0 {
		return errors.New("income must be greater than 0")
	}
	return nil
}

func validateDocuments(p []models.UploadDocument) error {
	if len(p) == 0 {
		return errors.New("no documents uploaded")
	}
	for i, doc := range p {
		if doc.Name == "" || doc.URL == "" {
			return fmt.Errorf("document %d is incomplete", i+1)
		}
	}
	return nil
}

func validateEducation(p []models.StudentAcademicQualification) error {

	if len(p) == 0 {
		return errors.New("no education history found")
	}
	for i, edu := range p {
		if edu.Degree == "" || edu.University == "" || edu.YearOfPassing == 0 {
			return fmt.Errorf("education history %d is incomplete", i+1)
		}
	}
	return nil
}

func validateAddresses(addresses []models.Address) error {
	if len(addresses) == 0 {
		return errors.New("no addresses provided")
	}

	hasPermanent := false
	hasCurrent := false

	for _, addr := range addresses {
		switch addr.Type {
		case "permanent":
			hasPermanent = true
		case "current":
			hasCurrent = true
		}

		if addr.Type == "" {
			return errors.New("address type is missing")
		}

		if addr.Street == "" || addr.City == "" || addr.State == "" || addr.Pincode == "" || addr.Country == "" {
			return fmt.Errorf("%s address is incomplete", addr.Type)
		}
	}

	if !hasPermanent {
		return errors.New("permanent address is required")
	}
	if !hasCurrent {
		return errors.New("current address is required")
	}

	return nil
}

// ============= Modify Application ==================
// UpsertStudentAddresses inserts or updates student addresses in the database.
func UpsertStudentAddresses(db *gorm.DB, studentID uint, addresses []models.AddressInput) error {
	//  Delete all existing addresses
	if err := db.Where("student_id = ?", studentID).Delete(&models.Address{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing addresses: %w", err)
	}

	var (
		permanentAdded bool
		currentAdded   bool
	)

	for _, addr := range addresses {
		if addr.Type != "permanent" && addr.Type != "current" {
			continue
		}

		// Skip fully empty address
		if addr.Street == "" && addr.City == "" && addr.State == "" && addr.Pincode == "" && addr.Country == "" {
			continue
		}

		// Enforce only one permanent and one current address
		if addr.Type == "permanent" && permanentAdded {
			return errors.New("more than one permanent addresses is provided")
		}
		if addr.Type == "current" && currentAdded {
			return errors.New("more than one current addresses is provided")
		}

		newAddr := models.Address{
			StudentID: studentID,
			Type:      addr.Type,
			Street:    addr.Street,
			City:      addr.City,
			State:     addr.State,
			Pincode:   addr.Pincode,
			Country:   addr.Country,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&newAddr).Error; err != nil {
			return fmt.Errorf("failed to create address: %w", err)
		}

		if addr.Type == "permanent" {
			permanentAdded = true
		} else if addr.Type == "current" {
			currentAdded = true
		}
	}

	return nil
}

// UpsertStudentDocuments inserts or updates student documents in the database.
func UpsertStudentDocuments(db *gorm.DB, studentID uint, documents []models.DocumentInput) error {
	//  Delete existing documents for the student
	if err := db.Where("student_id = ?", studentID).Delete(&models.UploadDocument{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing documents: %w", err)
	}

	//  Insert new documents
	for _, doc := range documents {
		if doc.Name == "" || doc.URL == "" {
			continue
		}

		newDoc := models.UploadDocument{
			StudentID: studentID,
			Name:      doc.Name,
			URL:       doc.URL,
			CreatedAt: time.Now(),
		}

		if err := db.Create(&newDoc).Error; err != nil {
			return fmt.Errorf("failed to create document: %w", err)
		}
	}

	return nil
}

// UpsertEducationHistory inserts or updates a student's education history in the database.
func UpsertEducationHistory(db *gorm.DB, studentID uint, history []models.EducationHistoryInput) error {
	// Delete existing education history for the student
	if err := db.Where("student_id = ?", studentID).Delete(&models.StudentAcademicQualification{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing education records: %w", err)
	}

	// Insert new records
	for _, edu := range history {
		// Skip empty entries
		if edu.Degree == "" && edu.University == "" && edu.Course == "" && edu.Grade == "" && edu.YearOfPassing == 0 {
			continue
		}

		newEdu := models.StudentAcademicQualification{
			StudentID:     studentID,
			Degree:        edu.Degree,
			University:    edu.University,
			YearOfPassing: edu.YearOfPassing,
			Grade:         edu.Grade,
			Course:        edu.Course,
			CreatedAt:     time.Now(),
		}
		if err := db.Create(&newEdu).Error; err != nil {
			return fmt.Errorf("failed to create education record: %w", err)
		}
	}

	return nil
}

func UpdateStudentProfileFields(studentProfile *models.StudentProfile, input models.StudentProfileInput) error {
	modify := false
	if input.FullName != "" {
		studentProfile.FullName = input.FullName
		modify = true
	}
	if input.Email != "" {
		studentProfile.Email = input.Email
		modify = true
	}
	if input.PhoneNumber != "" {
		studentProfile.PhoneNumber = input.PhoneNumber
		modify = true
	}
	if input.DateOfBirth != "" {
		// If DOB is provided, validate it
		if err := validateDOB(input.DateOfBirth); err != nil {

			return err
		}
		studentProfile.DateOfBirth = input.DateOfBirth
		modify = true
	}
	if input.Qualification != "" {
		studentProfile.Qualification = input.Qualification
		modify = true
	}
	if input.Category != "" {
		studentProfile.Category = input.Category
		modify = true
	}
	if input.Income > 0 {
		studentProfile.Income = input.Income
		modify = true
	}
	if input.Nationality != "" {
		studentProfile.Nationality = input.Nationality
		modify = true
	}
	if input.Gender != "" {
		studentProfile.Gender = input.Gender
		modify = true
	}
	if input.AadhaarNumber != "" {
		studentProfile.AadhaarNumber = input.AadhaarNumber
		modify = true
	}
	if modify {
		fmt.Println("true")
		studentProfile.CreatedAt = time.Now()
	}

	return nil
}

// ======================= search schemes ==================
// ApplySchemeFilters applies various filters to a database query for schemes.
func ApplySchemeFilters(query *gorm.DB, filter models.SchemeFilter) *gorm.DB {
	// Scheme Filters
	if filter.Name != nil {
		query = query.Where("schemes.name ILIKE ?", "%"+*filter.Name+"%")
	}
	if filter.Status != nil {
		query = query.Where("schemes.status = ?", *filter.Status)
	}
	if filter.MinAmount != nil {
		query = query.Where("schemes.amount >= ?", *filter.MinAmount)
	}
	if filter.MaxAmount != nil {
		query = query.Where("schemes.amount <= ?", *filter.MaxAmount)
	}
	if filter.StartAfter != nil {
		query = query.Where("schemes.start_date >= ?", *filter.StartAfter)
	}
	if filter.EndBefore != nil {
		query = query.Where("schemes.end_date <= ?", *filter.EndBefore)
	}

	// Eligibility Filters
	if filter.Gender != nil {
		query = query.Where("eligibilities.gender = ?", *filter.Gender)
	}
	if filter.AcademicQualification != nil {
		query = query.Where("eligibilities.academic_qualification = ?", *filter.AcademicQualification)
	}
	if filter.IncomeLimit != nil {
		query = query.Where("eligibilities.income_limit >= ?", *filter.IncomeLimit)
	}
	if filter.Category != nil {
		query = query.Where("eligibilities.category = ?", *filter.Category)
	}

	return query
}

// ====================== Pagination/url create ===================
func GetPagination(c *gin.Context) (models.PaginationInput, int64) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	return models.PaginationInput{
		Page:  page,
		Limit: limit,
	}, (page - 1) * limit
}

func BuildPaginationMeta(c *gin.Context, pagination models.PaginationInput, totalCount int64) *models.PaginationMeta {
	totalPages := int64(math.Ceil(float64(totalCount) / float64(pagination.Limit)))

	params := c.Request.URL.Query()
	basePath := c.Request.URL.Path

	var previous, next string
	if pagination.Page > 1 {
		previous = buildURL(basePath, params, pagination.Page-1)
	}
	if pagination.Page < totalPages {
		next = buildURL(basePath, params, pagination.Page+1)
	}

	return &models.PaginationMeta{
		ResourceCount: int(totalCount),
		TotalPages:    totalPages,
		Page:          pagination.Page,
		Limit:         pagination.Limit,
		Previous:      previous,
		Next:          next,
	}
}

func buildURL(basePath string, params url.Values, page int64) string {
	params.Set("page", strconv.FormatInt(page, 10))
	return basePath + "?" + params.Encode()
}

// =========================== validate function====================
func validateDOB(dobStr string) error {
	dob, err := time.Parse("02/01/2006", dobStr)
	if err != nil {
		return errors.New("date of birth is invalid, expected format DD/MM/YYYY")
	}
	if dob.After(time.Now()) {
		return errors.New("date of birth cannot be in the future")
	}
	return nil
}

// ============================ validator ===========================
type Validator interface {
	Name() string
	Validate(app *models.Application) error
}

type ValidationChain interface {
	AddValidator(validator Validator) ValidationChain
	Validate(app *models.Application) error
}

type validationChain struct {
	validators []Validator
}

func NewValidationChain() ValidationChain {
	return &validationChain{
		validators: make([]Validator, 0),
	}
}

func (vc *validationChain) AddValidator(validator Validator) ValidationChain {
	vc.validators = append(vc.validators, validator)
	return vc
}
