package oraaq

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pgvanniekerk/ezQue/api"
)

func NewEnqueuer(db *sql.DB, queueName string) *Enqueuer {
	return &Enqueuer{
		db:         db,
		queueName:  queueName,
		enqueueSql: enqueueSql,
	}
}

type Enqueuer struct {

	// db is a pointer to the SQL database connection. Note, this acts
	// as a connection pool by default, and is safe for concurrent use.
	db *sql.DB

	// The Oracle Advance Queue that the Enqueuer is bound to.
	queueName string

	enqueueSql string
}

// NewMessage returns a new instance of `Message` that implements the `api.Message` interface.
// It initializes the `ID` field with an empty byte slice and the `Content` field with an empty string.
// Clients can use the returned `Message` instance to set the ID and content as needed.
func (e *Enqueuer) NewMessage() api.Message[Message] {
	return &Message{}
}

// Enqueue enqueues a message to the Oracle Advanced Queue.
// It starts a new transaction, performs SQL to enqueue the message using the provided context and
// message content, and commits the transaction. If any error occurs during the process, it rolls back
// the transaction and returns an error.
func (e *Enqueuer) Enqueue(ctx context.Context, msg api.Message[Message]) error {

	// Start a new transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Perform SQL to enqueue message
	_, err = tx.ExecContext(ctx, e.enqueueSql, e.queueName, msg.Text())
	if err != nil {
		// Rollback transaction in case of an error
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("enqueue failed: %v, failed to rollback: %w", err, rollbackErr)
		}
		return fmt.Errorf("failed to enqueue message: %w", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (e *Enqueuer) Disconnect(_ context.Context) error {

	err := e.db.Close()
	if err != nil {
		return err
	}

	e.db = nil
	e.enqueueSql = ""
	e.queueName = ""
	return nil
}
