package oraaq

import (
	"fmt"
	"github.com/pgvanniekerk/ezQue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestConnector ensures that OracleAqJms can be provided
// to the ezQue.Connect function.
func TestConnector(t *testing.T) {
	_, _ = ezQue.Connect(OracleAqJms, nil)
}

// TestQueue ensures that Queue can be provided to the
// ezQue.Connect function as for the O type parameter(Options).
func TestQueue(t *testing.T) {
	_, _ = ezQue.Connect(OracleAqJms, Queue("testQueue"))
}

// TestAuthenticatedWith ensures that it correctly sets the config values
// and is a valid functional option to provide to the Queue method.
func TestAuthenticatedWith(t *testing.T) {

	// Test that the username and password values provided are set
	const (
		username = "testUser"
		password = "testPassword"
	)

	urlOpts := &urlOptions{}
	authFunc := AuthenticatedWith(username, password)
	authFunc(urlOpts)

	require.Equal(t, username, urlOpts.Username, "The username in urlOptions did not match the expected value")
	require.Equal(t, password, urlOpts.Password, "The password in urlOptions did not match the expected value")

	// Test that AuthenticatedWith can be provided as a functional option - by compiling, this passes
	_, _ = ezQue.Connect(OracleAqJms, Queue("testQueue", authFunc))
}

// TestLocatedAt ensures that it correctly sets the config values
func TestLocatedAt(t *testing.T) {
	const (
		server        = "localhost"
		port   uint16 = 8080
	)

	urlOpts := &urlOptions{}
	locFunc := LocatedAt(server, port)
	locFunc(urlOpts)

	require.Equal(t, server, urlOpts.Server, "The server in urlOptions did not match the expected value")
	require.Equal(t, port, urlOpts.Port, "The port number in urlOptions did not match the expected value")
	_, _ = ezQue.Connect(OracleAqJms, Queue("testQueue", locFunc))
}

// TestUsingService ensures that it correctly sets the config values
func TestUsingService(t *testing.T) {
	const service = "someService"

	urlOpts := &urlOptions{}
	servFunc := UsingService(service)
	servFunc(urlOpts)

	require.Equal(t, service, urlOpts.Service, "The service in urlOptions did not match the expected value")
	_, _ = ezQue.Connect(OracleAqJms, Queue("testQueue", servFunc))
}

// TestUsingSID ensures that it correctly sets the config values
func TestUsingSID(t *testing.T) {
	const sid = "someSID"

	urlOpts := &urlOptions{
		keyVals: make(map[string]string),
	}
	sidFunc := UsingSID(sid)
	sidFunc(urlOpts)

	require.Equal(t, sid, urlOpts.keyVals["SID"], "The SID in urlOptions did not match the expected value")
	_, _ = ezQue.Connect(OracleAqJms, Queue("testQueue", sidFunc))
}

// TestUsingJdbcString ensures that it correctly sets the config values
//func TestUsingJdbcString(t *testing.T) {
//	const jdbc = "someJDBC"
//
//	urlOpts := &urlOptions{
//		keyVals: make(map[string]string),
//	}
//	jdbcFunc := UsingJdbcString(jdbc)
//	jdbcFunc(urlOpts)
//
//	require.Equal(t, jdbc, urlOpts.keyVals["connStr"], "The JDBC string in urlOptions did not match the expected value")
//	_, err := ezQue.Connect(OracleAqJms, Queue("testQueue", jdbcFunc))
//	if err != nil {
//		t.Fatal(err)
//	}
//}

// TestWithURLOptions ensures that it correctly sets the config values
func TestWithURLOptions(t *testing.T) {
	urlOpts := &urlOptions{
		keyVals: make(map[string]string),
	}
	urlOptions := map[string]string{
		"opt1": "val1",
		"opt2": "val2",
	}
	urlOptFunc := WithURLOptions(urlOptions)
	urlOptFunc(urlOpts)

	for key, val := range urlOptions {
		require.Equal(t, val, urlOpts.keyVals[key], fmt.Sprintf("The value of %s in urlOptions did not match the expected value", key))
	}
	_, _ = ezQue.Connect(OracleAqJms, Queue("testQueue", urlOptFunc))
}

//
// connect
//

type ConnectTestSuite struct {
	suite.Suite
}

func TestConnectTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectTestSuite))
}

func (suite *ConnectTestSuite) Test_FailOnNilOptions() {
	_, _, err := connect(nil)
	assert.NotNil(suite.T(), err)
}

func (suite *ConnectTestSuite) Test_FailOnEmptyQueueName() {
	_, _, err := connect(func() OptionFunc {
		return func() Options {
			return Options{
				queueName: "",
				urlOpts: []UrlOptionFunc{
					func(options *urlOptions) {
						options.keyVals = make(map[string]string)
						options.keyVals["DUMMYKEY"] = "DUMMYVAL"
					},
				},
			}
		}
	}())
	assert.NotNil(suite.T(), err)
}

// Test_FailOnErrorWhenSqlOpen tests that an error is returned if sql.Open returns an error value.
func (suite *ConnectTestSuite) Test_ErrorOnSqlOpenFailure() {
	_, _, err := connect(func() OptionFunc {
		return func() Options {
			return Options{
				queueName: "VALID_QUEUE_NAME",
				urlOpts: []UrlOptionFunc{
					func(options *urlOptions) {
						options.keyVals = make(map[string]string)
						options.keyVals["DUMMYKEY"] = "DUMMYVAL"
					},
				},
			}
		}
	}())
	assert.NotNil(suite.T(), err)
}

//
// connectEnqueue
//

func TestConnectEnqueueTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectEnqueueTestSuite))
}

type ConnectEnqueueTestSuite struct {
	suite.Suite
}

// Test_FailOnNilOptions tests that an error is returned if no connectOptions have been provided
func (suite *ConnectEnqueueTestSuite) Test_FailOnNilOptions() {
	_, err := connectEnqueue(nil)
	assert.NotNil(suite.T(), err)
}

// Test_FailOnEmptyQueueName tests that an error is returned if an empty string is provided
// as the queueName. Add "dummy" values as urlOptions, to ensure failure is only due to
// empty queueName.
func (suite *ConnectEnqueueTestSuite) Test_FailOnEmptyQueueName() {
	_, err := connectEnqueue(func() OptionFunc {
		return func() Options {
			return Options{
				queueName: "",
				urlOpts: []UrlOptionFunc{
					func(options *urlOptions) {
						options.keyVals = make(map[string]string)
						options.keyVals["DUMMYKEY"] = "DUMMYVAL"
					},
				},
			}
		}
	}())
	assert.NotNil(suite.T(), err)
}

// Test_FailOnErrorWhenSqlOpen tests that an error is returned if sql.Open returns an error value.
func (suite *ConnectEnqueueTestSuite) Test_ErrorOnSqlOpenFailure() {
	_, err := connectEnqueue(func() OptionFunc {
		return func() Options {
			return Options{
				queueName: "VALID_QUEUE_NAME",
				urlOpts: []UrlOptionFunc{
					func(options *urlOptions) {
						options.keyVals = make(map[string]string)
						options.keyVals["DUMMYKEY"] = "DUMMYVAL"
					},
				},
			}
		}
	}())
	assert.NotNil(suite.T(), err)
}

//
// connectDequeue
//

func TestConnectDequeueTestSuite(t *testing.T) {
	suite.Run(t, new(ConnectDequeueTestSuite))
}

type ConnectDequeueTestSuite struct {
	suite.Suite
}

// Test_FailOnNilOptions tests that an error is returned if no connectOptions have been provided
func (suite *ConnectDequeueTestSuite) Test_FailOnNilOptions() {
	_, err := connectDequeue(nil)
	assert.NotNil(suite.T(), err)
}

// Test_FailOnEmptyQueueName tests that an error is returned if an empty string is provided
// as the queueName. Add "dummy" values as urlOptions, to ensure failure is only due to
// empty queueName.
func (suite *ConnectDequeueTestSuite) Test_FailOnEmptyQueueName() {
	_, err := connectDequeue(func() OptionFunc {
		return func() Options {
			return Options{
				queueName: "",
				urlOpts: []UrlOptionFunc{
					func(options *urlOptions) {
						options.keyVals = make(map[string]string)
						options.keyVals["DUMMYKEY"] = "DUMMYVAL"
					},
				},
			}
		}
	}())
	assert.NotNil(suite.T(), err)
}

// Test_ErrorOnSqlOpenFailure tests that an error is returned if sql.Open returns an error value.
func (suite *ConnectDequeueTestSuite) Test_ErrorOnSqlOpenFailure() {
	_, err := connectDequeue(func() OptionFunc {
		return func() Options {
			return Options{
				queueName: "VALID_QUEUE_NAME",
				urlOpts: []UrlOptionFunc{
					func(options *urlOptions) {
						options.keyVals = make(map[string]string)
						options.keyVals["DUMMYKEY"] = "DUMMYVAL"
					},
				},
			}
		}
	}())
	assert.NotNil(suite.T(), err)
}
