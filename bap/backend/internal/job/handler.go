package job

// import (
// 	"github.com/sirupsen/logrus"

// 	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/clients"
// 	jobDb "github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/database/mongodb/job"
// 	"github.com/ONEST-Network/Whatsapp-Chatbot/bap/backend/pkg/utils/random"
// )

// type Interface interface {
// 	CreateJob(payload *jobDb.Job) error
// }

// type Job struct {
// 	clients *clients.Clients
// }

// func NewJob(clients *clients.Clients) Interface {
// 	return &Job{
// 		clients: clients,
// 	}
// }

// func (j *Job) CreateJob(payload *jobDb.Job) error {
// 	logrus.Infof("Received request to create a new job for business: %s", payload.BusinessID)

// 	// populate job id
// 	payload.ID = random.GetRandomString(7)

// 	if err := j.clients.JobClient.CreateJob(payload); err != nil {
// 		return err
// 	}

// 	return nil
// }
