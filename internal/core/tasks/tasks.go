// Package tasks provides the core contracts for background job definitions.
package tasks

import "github.com/jalevin/gottl/internal/core/mailer"

type TaskID string

var (
	TaskIDSendEmail         TaskID = "send_email"
	TaskIDDeleteExpiredData TaskID = "delete_expired_data"
)

type Task struct {
	ID      TaskID
	Payload any
}

type Queue interface {
	Enqueue(task Task) error
}

func NewEmailTask(message mailer.Message) Task {
	return Task{
		ID:      TaskIDSendEmail,
		Payload: message,
	}
}

func NewDeleteExpiredDataTask() Task {
	return Task{
		ID:      TaskIDDeleteExpiredData,
		Payload: nil,
	}
}
