package job

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/job"
)

type CreateJobRequest struct {
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Type           job.JobType     `json:"type"`
	Vacancies      int             `json:"vacancies"`
	SalaryRange    job.SalaryRange `json:"salaryRange"`
	ApplicationIDs []string        `json:"applicationIds"`
	WorkHours      job.WorkHours   `json:"workHours"`
	WorkDays       job.WorkDays    `json:"workDays"`
	Eligibility    job.Eligibility `json:"eligibility"`
	Location       job.Location    `json:"location"`
	BusinessID     string          `json:"businessId"`
}