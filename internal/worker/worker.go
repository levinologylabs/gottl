// Package worker provides the worker that processes the tasks.
package worker

import (
	"context"

	"github.com/jalevin/gottl/internal/core/mailer"
	"github.com/jalevin/gottl/internal/core/tasks"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/rs/zerolog"
)

type worker struct {
	l      zerolog.Logger
	sender mailer.Sender
	db     *db.QueriesExt
}

func newWorker(l zerolog.Logger, db *db.QueriesExt, sender mailer.Sender) *worker {
	return &worker{
		l:      l,
		db:     db,
		sender: sender,
	}
}

func (w *worker) run(ctx context.Context, ch <-chan tasks.Task) {
	defer func() {
		if r := recover(); r != nil {
			go w.run(ctx, ch)
		}
	}()

	select {
	case <-ctx.Done():
		return
	case task := <-ch:
		switch task.ID {
		case tasks.TaskIDSendEmail:
			w.sendEmail(ctx, task)
		case tasks.TaskIDDeleteExpiredData:
			w.deleteExpiredData(ctx, task)
		}
	}
}

func (w *worker) sendEmail(ctx context.Context, task tasks.Task) {
	data, ok := task.Payload.(mailer.Message)
	if !ok {
		w.l.Error().Ctx(ctx).Msg("invalid payload for email task")
	}

	err := w.sender.Send(data)
	if err != nil {
		w.l.Error().Ctx(ctx).Err(err).
			Str("email", data.To).
			Str("subject", data.Subject).
			Msg("failed to send email")
	}
}

func (w *worker) deleteExpiredData(ctx context.Context, task tasks.Task) {
	panic("not implemented")
}
