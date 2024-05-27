package oraaq

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	go_ora "github.com/sijms/go-ora/v2"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestDequeuerTestSuite(t *testing.T) {
	suite.Run(t, new(DequeuerTestSuite))
}

type DequeuerTestSuite struct {
	suite.Suite
	db        *sql.DB
	queueName string
}

func (suite *DequeuerTestSuite) SetupSuite() {

	// How to retrieve connection details from somewhere
	server := os.Getenv("DB_SERVER")
	portString := os.Getenv("DB_PORT")
	port, err := strconv.Atoi(portString)
	if err != nil {
		suite.T().Fatalf("Invalid port number: %v", err)
	}
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	sid := os.Getenv("DB_SID")
	urlOptions := map[string]string{
		"SID": sid,
	}

	// Connect to test oracle database
	connStr := go_ora.BuildUrl(server, port, "", username, password, urlOptions)
	db, err := sql.Open("oracle", connStr)
	if err != nil {
		suite.T().Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		suite.T().Errorf("Error connecting to database: %v", err)
	}
	suite.db = db

	// Execute the setupTestQueue.sql script
	// Read setup script
	setupScript, err := os.ReadFile(filepath.Join(".", "setupTestQueue.sql"))
	if err != nil {
		suite.T().Fatalf("Error reading setup script: %v", err)
	}

	// Split the file contents into individual SQL statements
	sqlStatements := strings.Split(string(setupScript), "/")

	// Execute each SQL statement
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			suite.T().Errorf("Error executing statement: %s\nError: %s\n", stmt, err)
		}
	}

}

func (suite *DequeuerTestSuite) TearDownSuite() {

	// Read teardown script
	teardownScript, err := os.ReadFile(filepath.Join(".", "tearDownTestQueue.sql"))
	if err != nil {
		suite.T().Fatalf("Error reading teardown script: %v", err)
	}

	// Split the file contents into individual SQL statements
	sqlStatements := strings.Split(string(teardownScript), "/")

	// Execute each SQL statement
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := suite.db.Exec(stmt)
		if err != nil {
			suite.T().Errorf("Error executing statement: %s\nError: %s\n", stmt, err)
		}
	}

	// Close the database connection
	err = suite.db.Close()
	if err != nil {
		suite.T().Fatalf(err.Error())
	}

}

// SetupTest is run before each test in the suite
func (suite *DequeuerTestSuite) SetupTest() {
	// Prepare the SQL statement to purge the queue
	purgeQueueSQL := `
Delete
From text_msg_queue_table
	`

	// Execute the SQL statement
	_, err := suite.db.Exec(purgeQueueSQL)
	if err != nil {
		suite.T().Fatalf("Failed to purge the queue table: %v", err)
	}
}

func (suite *DequeuerTestSuite) TestNewDequeuer() {
	db := &sql.DB{} // usually you'd use sql.Open to get a *sql.DB
	const queueName = "testQueue"

	dequeuer := NewDequeuer(db, queueName)

	// check if the db and queueName properties of the dequeuer match the parameters that you passed
	require.Equal(suite.T(), db, dequeuer.db, "The dequeuer's db property does not match the expected db")
	require.Equal(suite.T(), queueName, dequeuer.queueName, "The dequeuer's queueName property does not match the expected queue name")
}

func (suite *DequeuerTestSuite) TestDequeue_BeginTxError() {
	const queueName = "testQueue"

	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(suite.T(), err, "An error was not expected when opening a stub database connection")

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	// Create a Dequeuer and call Dequeue
	dequeuer := NewDequeuer(db, queueName)
	_, err = dequeuer.Dequeue(context.Background())

	require.Error(suite.T(), err, "An error was expected when calling Dequeue because BeginTx should fail")
}

func (suite *DequeuerTestSuite) TestDequeue_QueryContextError() {
	const queueName = "testQueue"

	// Create a new mock database
	db, mock, err := sqlmock.New()
	require.NoError(suite.T(), err, "An error was not expected when opening a stub database connection")

	mock.ExpectBegin()
	mock.ExpectQuery(dequeueSQL).WithArgs(queueName).WillReturnError(sql.ErrConnDone)

	// Create a Dequeuer and call Dequeue
	dequeuer := NewDequeuer(db, queueName)
	_, err = dequeuer.Dequeue(context.Background())

	require.Error(suite.T(), err, "An error was expected when calling Dequeue because QueryContext should fail")
}

func (suite *DequeuerTestSuite) TestDequeueEmptyQueue() {
	// create an instance of Dequeuer
	dequeuer := NewDequeuer(suite.db, "text_msg_queue")

	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Cancel context when operation finishes or timeout occurs

	// Attempt to Dequeue from an empty queue
	deqMsg, err := dequeuer.Dequeue(ctx)

	// The operation should return an error
	suite.Error(err, "Expected an error when dequeuing from an empty queue")
	// Assert that the error is a deadline exceeded error indicating a timeout
	suite.EqualError(err, context.DeadlineExceeded.Error())

	if err != nil {
		// If an error occurred, the returned DequeueMessage should be nil
		suite.Nil(deqMsg, "DequeueMessage should be nil when there is an error")
	}
}

func (suite *DequeuerTestSuite) TestDequeueWithMessageAndRollback() {
	// Create an instance of Dequeuer
	dequeuer := NewDequeuer(suite.db, "text_msg_queue")

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Cancel context when operation finishes or timeout occurs

	// Enqueue a message directly using PL/SQL
	enqueueSQL := `
		BEGIN
		   -- Declare variables
		   DECLARE
			  v_enqueue_options  DBMS_AQ.ENQUEUE_OPTIONS_T;
			  v_message_properties DBMS_AQ.MESSAGE_PROPERTIES_T;
			  v_message SYS.AQ$_JMS_TEXT_MESSAGE;
			  v_msgid RAW(16);
		   BEGIN
			  -- Create a text message
			  -- Initialize JMS text message type
			  v_message := SYS.AQ$_JMS_TEXT_MESSAGE.construct;
			  -- Set the text payload of the message
			  v_message.set_text('test message');
			  -- Enqueue message into the queue
			  DBMS_AQ.ENQUEUE(
				 queue_name          => 'text_msg_queue',
				 enqueue_options     => v_enqueue_options,
				 message_properties  => v_message_properties,
				 payload             => v_message,
				 msgid               => v_msgid
			  );
			  
			  COMMIT;
		   END;
		END;
	`

	_, err := suite.db.Exec(enqueueSQL)
	suite.NoError(err, "Failed to enqueue message")

	// Dequeue the message
	deqMsg, err := dequeuer.Dequeue(ctx)
	suite.NoError(err, "Failed to dequeue message")
	suite.NotNil(deqMsg, "Dequeued message should not be nil")

	// Assert received message
	msg := deqMsg.Message().Raw()
	suite.Equal("test message", msg.Text(), "Dequeued message should equal enqueued message")

	// NAck/rollback the message
	err = deqMsg.NAck(ctx)
	suite.NoError(err, "Failed to NAck message")

	// Re-dequeue the message
	deqMsg2, err := dequeuer.Dequeue(ctx)
	suite.NoError(err, "Failed to re-dequeue message")
	suite.NotNil(deqMsg2, "Re-dequeued message should not be nil")

	// Assert re-dequeued message
	msg2 := deqMsg2.Message().Raw()
	suite.Equal("test message", msg2.Text(), "Re-dequeued message should still be identical")
}

func (suite *DequeuerTestSuite) TestDequeueWithMessageAndAck() {
	// Create an instance of Dequeuer
	dequeuer := NewDequeuer(suite.db, "text_msg_queue")

	// Create a context with a timeout of 10 seconds
	enqCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Enqueue a message directly using PL/SQL
	enqueueSQL := `
	BEGIN
		DECLARE
			v_enqueue_options  DBMS_AQ.ENQUEUE_OPTIONS_T;
			v_message_properties DBMS_AQ.MESSAGE_PROPERTIES_T;
			v_message SYS.AQ$_JMS_TEXT_MESSAGE;
			v_msgid RAW(16);
		BEGIN
			-- Create a text message
			-- Initialize JMS text message type
			v_message := SYS.AQ$_JMS_TEXT_MESSAGE.construct;
			-- Set the text payload of the message
			v_message.set_text('test message');

			-- Enqueue message into the queue
			DBMS_AQ.ENQUEUE(
				queue_name          => 'text_msg_queue',
				enqueue_options     => v_enqueue_options,
				message_properties  => v_message_properties,
				payload             => v_message,
				msgid               => v_msgid
			);

			COMMIT;
		END;
	END;
	`

	_, err := suite.db.ExecContext(enqCtx, enqueueSQL)
	suite.NoError(err, "Failed to enqueue message")

	// Dequeue the message
	deqCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	deqMsg, err := dequeuer.Dequeue(deqCtx)
	suite.NoError(err, "Failed to dequeue message")
	suite.NotNil(deqMsg, "Dequeued message should not be nil")

	// Assert received message
	receivedMsg := deqMsg.Message().Text()
	suite.Equal("test message", receivedMsg, "Dequeued message should equal enqueued message")

	// Ack/commit the message
	err = deqMsg.Ack(deqCtx)
	suite.NoError(err, "Failed to Ack message")

	// Try to re-dequeue the message. This should fail with a context deadline exceeded error if the Ack worked.
	reCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // always cancel context when done

	newDeqMsg, err := dequeuer.Dequeue(reCtx)

	if newDeqMsg != nil {
		suite.T().Errorf("newDeqMsg should be nil: %v", newDeqMsg)
	}

	// We expect a timeout error because the message has been acked and should no longer be in the queue.
	suite.Error(err, "Expected error when trying to re-dequeue the message")
	suite.Equal(err, context.DeadlineExceeded)
}
