package ezQue

import (
	"github.com/pgvanniekerk/ezQue/api"
)

// queueConnector is a type alias for a function that constructs and initializes Enqueuer and Dequeuer objects.
// It takes an options parameter of any type O which is used for configuration.
// This function is expected to return an instance of Enqueuer[M], Dequeuer[M], and an error if there was an issue with construction.
// It is intended to be used as a factory function for creating specific implementations of the Enqueuer/Dequeuer.
type queueConnector[R, O any] func(options O) (api.Enqueuer[R], api.Dequeuer[R], error)

// Connect is a generic function that establishes a connection to a queue given a function
// for creating specific Enqueuer and Dequeuer implementations (queueConnector) and
// options for the Enqueuer/Dequeuer objects.
//
// The function takes a queueConnector of type qc that is responsible for creating the Enqueuer and Dequeuer
// objects used by the Queue interface. It also takes an options parameter of type O
// that is used to configure the Enqueuer and Dequeuer objects.
//
// The Connect function will call the queueConnector function with the provided options. If the queueConnector function
// returns an error, Connect will propagate this error and return.
//
// If no error occurred during the creation of the Enqueuer and Dequeuer, Connect will create a new Queue
// implementation using the provided Enqueuer and Dequeuer. It then returns the newly created Queue and a nil value
// for the error.
//
// Callers should use this function when they want to set up a Queue with specific Enqueuer and Dequeuer implementations.
func Connect[R, O any](qc queueConnector[R, O], options O) (Queue[R], error) {

	enqueuer, dequeuer, err := qc(options)
	if err != nil {
		return nil, err
	}

	q := &queue[R]{
		enqueuer: enqueuer,
		dequeuer: dequeuer,
	}

	return q, nil
}
