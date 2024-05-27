package oraaq

import (
	"context"
	"database/sql"
	go_ora "github.com/sijms/go-ora/v2"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

type EnqueuerTestSuite struct {
	suite.Suite
	db        *sql.DB
	queueName string
}

func TestEnqueuerTestSuite(t *testing.T) {
	suite.Run(t, new(EnqueuerTestSuite))
}

func (suite *EnqueuerTestSuite) SetupSuite() {
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

	suite.queueName = "text_msg_queue"
}

func (suite *EnqueuerTestSuite) TearDownSuite() {

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

func (suite *EnqueuerTestSuite) TestEnqueue() {
	enqueuer := NewEnqueuer(suite.db, suite.queueName)

	message := &Message{Content: "test message"}
	// Create a context with a timeout of 5 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // Cancel context when operation finishes or timeout occurs

	err := enqueuer.Enqueue(ctx, message)
	suite.NoError(err, "Failed to enqueue message")

	// Begin a new transaction for dequeue
	tx, err := suite.db.BeginTx(ctx, nil)
	if err != nil {
		suite.T().Fatal(err)
	}

	var content go_ora.Clob
	var msgID string
	var errMsg sql.NullString

	// Execute the dequeue PL/SQL anonymous block
	_, err = tx.ExecContext(ctx, dequeueSQL, suite.queueName,
		go_ora.Out{Dest: &content, Size: 300000},
		go_ora.Out{Dest: &msgID, Size: 32},
		go_ora.Out{Dest: &errMsg, Size: 4000},
	)

	if err != nil {

		if strings.Contains(err.Error(), "ORA-01013") {
			// Wrap it as a context deadline exceeded
			err = context.DeadlineExceeded
		}

		suite.T().Fatal(err)
	}

	if errMsg.Valid {
		suite.T().Fatal(errMsg.String)
	}

	// Compare the content of the dequeued message with the original message
	expectedMessage := "test message"
	suite.Equal(expectedMessage, content.String, "The content of the dequeued message does not match the original message.")
}

func (suite *EnqueuerTestSuite) TestEnqueue_Error() {
	enqueuer := NewEnqueuer(suite.db, "pfft")

	// create a message with text
	message := &Message{Content: "test message"}
	ctx := context.Background()

	// enqueue a message into list queue -> should return an error
	err := enqueuer.Enqueue(ctx, message)
	suite.Error(err, "Expected error during enqueue because the queue does not exist (Enqueue should fail).")
}

func (suite *EnqueuerTestSuite) TestEnqueue_disconnect() {
	enqueuer := NewEnqueuer(suite.db, suite.queueName)

	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	message := &Message{Content: "test message"}

	err := enqueuer.Enqueue(ctx, message)

	// Expect error due to context being cancelled
	suite.Error(err, "Expected an error when attempting to enqueue with a cancelled context")
}
