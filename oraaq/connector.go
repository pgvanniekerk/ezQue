package oraaq

import (
	"database/sql"
	"fmt"
	"github.com/pgvanniekerk/ezQue/api"
	"github.com/pgvanniekerk/ezQue/internal/oraaq"
	go_ora "github.com/sijms/go-ora/v2"
)

// OracleAqJms is provided as a queueConnector to ezQueue.Connect method, to connect to an Oracle Advance Queue.
func OracleAqJms(options OptionFunc) (api.Enqueuer[oraaq.Message], api.Dequeuer[oraaq.Message], error) {
	return connect(options)
}

// connect establishes both enqueue and dequeue connections.
func connect(opts OptionFunc) (api.Enqueuer[oraaq.Message], api.Dequeuer[oraaq.Message], error) {

	// Initialise Enqueuer
	enq, err := connectEnqueue(opts)
	if err != nil {
		return nil, nil, err
	}

	// Initialise Dequeuer
	deq, err := connectDequeue(opts)
	if err != nil {
		return nil, nil, err
	}

	return enq, deq, nil
}

// connectEnqueue establishes an enqueue connection.
func connectEnqueue(options OptionFunc) (api.Enqueuer[oraaq.Message], error) {

	// Get the Options
	if options == nil {
		return nil, fmt.Errorf("oraaq: options is nil")
	}
	opts := options()

	// Validate queueName
	if opts.queueName == "" {
		return nil, fmt.Errorf("oraaq: queueName is empty")
	}

	// Build db URL
	urlOpts := &urlOptions{
		keyVals: make(map[string]string),
	}
	for _, opt := range opts.urlOpts {
		opt(urlOpts)
	}

	// connect to db
	db, err := sql.Open("oracle", go_ora.BuildUrl(urlOpts.Server, int(urlOpts.Port), urlOpts.Service, urlOpts.Username, urlOpts.Password, urlOpts.keyVals))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	enq := oraaq.NewEnqueuer(db, opts.queueName)
	return enq, nil
}

// connectDequeue establishes a dequeue connection.
func connectDequeue(options OptionFunc) (api.Dequeuer[oraaq.Message], error) {

	// Get the Options
	if options == nil {
		return nil, fmt.Errorf("oraaq: options is nil")
	}
	opts := options()

	// Validate queueName
	if opts.queueName == "" {
		return nil, fmt.Errorf("oraaq: queueName is empty")
	}

	// Build db URL
	urlOpts := &urlOptions{
		keyVals: make(map[string]string),
	}
	for _, opt := range opts.urlOpts {
		opt(urlOpts)
	}

	// connect to db
	db, err := sql.Open("oracle", go_ora.BuildUrl(urlOpts.Server, int(urlOpts.Port), urlOpts.Service, urlOpts.Username, urlOpts.Password, urlOpts.keyVals))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	deq := oraaq.NewDequeuer(db, opts.queueName)
	return deq, nil
}

// Options struct holds url options and queue name.
type Options struct {
	urlOpts   []UrlOptionFunc
	queueName string
	db        *sql.DB
}

// OptionFunc is a function type that returns Options.
type OptionFunc func() Options

// Queue returns an OptionFunc that holds url options and a queue name.
func Queue(queue string, urlOpts ...UrlOptionFunc) OptionFunc {
	return func() Options {
		return Options{
			urlOpts:   urlOpts,
			queueName: queue,
		}
	}
}

// UrlOptionFunc is a function type to set urlOptions.
type UrlOptionFunc func(*urlOptions)

// urlOptions struct holds the URL information required for creating connections.
type urlOptions struct {
	Username string
	Password string
	Server   string
	Port     uint16
	Service  string
	keyVals  map[string]string
}

// AuthenticatedWith sets username and password for UrlOptionFunc.
func AuthenticatedWith(username string, password string) UrlOptionFunc {
	return func(opts *urlOptions) {
		opts.Username = username
		opts.Password = password
	}
}

// LocatedAt sets server and port for UrlOptionFunc.
func LocatedAt(server string, port uint16) UrlOptionFunc {
	return func(opts *urlOptions) {
		opts.Port = port
		opts.Server = server
	}
}

// UsingService sets service name for UrlOptionFunc.
func UsingService(service string) UrlOptionFunc {
	return func(opts *urlOptions) {
		opts.Service = service
	}
}

// UsingSID appends the specified SID to urlOptions keyVals.
func UsingSID(sid string) UrlOptionFunc {
	return func(opts *urlOptions) {
		opts.keyVals["SID"] = sid
	}
}

// UsingJdbcString appends the JDBC string to urlOptions keyVals.
//func UsingJdbcString(jdbc string) UrlOptionFunc {
//	return func(opts *urlOptions) {
//		opts.keyVals["connStr"] = jdbc
//	}
//}

// WithURLOptions appends specified URL options to urlOptions keyVals.
func WithURLOptions(urlOpts map[string]string) UrlOptionFunc {
	return func(opts *urlOptions) {
		for key, val := range urlOpts {
			opts.keyVals[key] = val
		}
	}
}
