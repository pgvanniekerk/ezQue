package oraaq

import (
	"context"
	"database/sql"
	"github.com/pgvanniekerk/ezQue/api"
)

type DequeueMessage struct {
	message Message
	tx      *sql.Tx
}

func (d *DequeueMessage) Message() api.Message[Message] {
	return &d.message
}

func (d *DequeueMessage) Ack(_ context.Context) error {
	return d.tx.Commit()
}

func (d *DequeueMessage) NAck(_ context.Context) error {
	return d.tx.Rollback()
}
