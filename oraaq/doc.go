// Package oraaq provides Go functions and types for working with Oracle Advance Queues
// (AQ). This includes operations such as enqueuing and dequeuing messages.
//
// Central to the package is the OracleAqJms function. This function takes in an OptionFunc parameter
// for customizing the connection options and returns Enqueuer and Dequeuer instances associated
// with the specified Oracle AQ.
//
// To configure various Oracle AQ options, the package defines an OptionFunc type. This type is a
// function that returns configured Options. Also provided are several helper function types such as
// UrlOptionFunc to set specific options.
//
// The package also provides low-level functions, connectEnqueue and connectDequeue, to initiate
// enqueue and dequeue connections individually.
//
// Note: This package relies on other packages, namely "ezQue/api", "ezQue/internal/oraaq",
// and "github.com/sijms/go-ora/v2". It must be used in the context where these packages are accessible.
package oraaq
