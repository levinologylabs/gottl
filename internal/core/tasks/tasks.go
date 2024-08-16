// Package tasks provides the core contracts for background job definitions.
package tasks

type TaskID string

var (
	TaskIDSendEmail         TaskID = "send_email"
	TaskIDDeleteExpiredData TaskID = "delete_expired_data"
)

type TaskDataSendEmail struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Task struct {
	ID      TaskID
	Payload any
}

type Queue interface {
	Enqueue(task Task) error
}
