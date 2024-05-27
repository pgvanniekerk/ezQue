package oraaq

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/pgvanniekerk/ezQue/api"
	go_ora "github.com/sijms/go-ora/v2"
	"strings"
)

func NewDequeuer(db *sql.DB, queueName string) *Dequeuer {
	return &Dequeuer{
		db:         db,
		queueName:  queueName,
		dequeueSql: dequeueSQL,
	}
}

// Dequeuer represents a type that dequeues messages from an Oracle Advanced Queue.
// It contains a pointer to the SQL database connection and the name of the queue it is bound to.
type Dequeuer struct {

	// db is a pointer to the SQL database connection. Note, this acts
	// as a connection pool by default, and is safe for concurrent use.
	db *sql.DB

	// The Oracle Advance Queue that the Dequeuer is bound to.
	queueName string

	dequeueSql string
}

// Dequeue retrieves a message from the Oracle Advanced Queue using a given transaction.
// It waits indefinitely until a message is returned or until the context is cancelled.
// The method begins a new transaction, executes the dequeue PL/SQL anonymous block, and
// reads the message data from the result set. It then builds a DequeueMessage object with
// the message data and the transaction. The result set is closed before returning the
// DequeueMessage object.
func (d *Dequeuer) Dequeue(ctx context.Context) (api.DequeueMessage[Message], error) {

	// Begin a new transaction that will be passed into the
	// DequeueMessage object to allow Commit/Rollback.
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var content go_ora.Clob
	var msgID string
	var errMsg sql.NullString

	// Execute the dequeue PL/SQL anonymous block, waiting
	// forever until a message is returned or until the context
	// has been cancelled.
	_, err = tx.ExecContext(ctx, d.dequeueSql,
		d.queueName,
		go_ora.Out{Dest: &content, Size: 300000},
		go_ora.Out{Dest: &msgID, Size: 32},
		go_ora.Out{Dest: &errMsg, Size: 4000},
	)
	if err != nil {

		// Check if the error is a 'cancel of current operation' from Oracle
		if strings.Contains(err.Error(), "ORA-01013") {
			// Wrap it as a context deadline exceeded
			err = context.DeadlineExceeded
		}

		return nil, err
	} else if (errMsg != sql.NullString{}) {
		return nil, fmt.Errorf("error occurred during dequeue: %s", errMsg.String)
	}

	// Decode hex string to byte slice
	decoded, err := hex.DecodeString(msgID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode msgID: %w", err)
	}

	// Read the message data from the result set
	var msgIDArray [16]byte
	copy(msgIDArray[:], decoded)
	message := Message{
		ID:      msgIDArray,
		Content: content.String,
	}

	// Build DequeueMessage
	deqMsg := &DequeueMessage{
		message: message,
		tx:      tx,
	}

	return deqMsg, nil
}

func (d *Dequeuer) Disconnect(_ context.Context) error {

	err := d.db.Close()
	if err != nil {
		return err
	}

	d.db = nil
	d.dequeueSql = ""
	d.queueName = ""
	return nil
}
