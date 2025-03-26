package business

import (
	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/workerProfile"
)

type AddWorkerProfileRequest struct {
	ID                 string                        `json:"id"`
	Name               string                        `json:"name"`
	Phone              string                        `json:"phone"`
	Email              string                        `json:"email"`
	Age                int                           `json:"age"`
	Gender             workerProfile.Gender          `json:"gender"`
	PreferredLanguages []workerProfile.Language      `json:"preferred_languages"`
	PreferredJobRoles  []workerProfile.JobRole       `json:"preferred_job_roles"`
	Experience         int                           `json:"experience"`
	Skills             []string                      `json:"skills"`
	Certifications     []workerProfile.Certification `json:"certifications"`
	Credentials        []workerProfile.Credential    `json:"credentials"`
	Location           workerProfile.Location        `json:"location"`
}
