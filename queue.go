package ezQue

import (
	"context"
	"github.com/pgvanniekerk/ezQue/api"
	"golang.org/x/sync/errgroup"
)

// Queue represents a generic queue interface where M represents any type.
// The Queue interface provides three methods for interaction:
// Enqueue for adding elements, Dequeue for retrieving elements, and Disconnect for terminating the connection.
type Queue[R any] interface {

	// NewMessage returns a new instance of `api.Message[R]` interface. It can be used to create
	// messages of type `R` for enqueuing into the queue. The returned message implements the
	// `Raw` and `Text` methods for retrieving the raw message and its text representation,
	// as well as the `SetRaw` and `SetText` methods for setting the raw message and its text.
	// The `NewMessage` method is part of the `Queue[R]` interface.
	NewMessage() api.Message[R]

	// Enqueue adds a new element of type M to the queue.
	// It takes in a context for handling cancellations and timeouts, and an element of type M to add to the queue.
	// It returns an error if the enqueuing operation fails.
	Enqueue(ctx context.Context, msg api.Message[R]) error

	// Dequeue retrieves and removes an element of type M from the queue.
	// Dequeue wraps the retrieved element in a DequeueMessage, which provides methods for
	// acknowledging the successful processing of the message (Ack) or negating its acknowledgment (NAck).
	// It blocks until a message is available or the provided context is cancelled or times out.
	// If the context is cancelled or times out, Dequeue will return a context cancellation error.
	// If the dequeuing operation fails for other reasons, it will return an error.
	// It's important to call Ack or NAck on the DequeMessage after processing it
	// to ensure the queue properly manages the message lifecycle.
	Dequeue(ctx context.Context) (api.DequeueMessage[R], error)

	// Disconnect closes the connection with the queue based on the provided context.
	// It should be called when the queue operations are no longer required.
	// It returns an error if there was an issue during the disconnection process.
	Disconnect(ctx context.Context) error
}

// queue is a generic data type that implements the Queue interface.
// It provides enqueue and dequeue operations for any type `M`.
// It's a concrete local implementation that delegates operations
// to its internal enqueuer and dequeuer.
type queue[R any] struct {
	enqueuer api.Enqueuer[R]
	dequeuer api.Dequeuer[R]
}

// NewMessage returns a new message object of type api.Message[R], created by the queue's enqueuer.
// The returned message object can be used to set the message's raw data and text.
// This method does not enqueue the message to the queue.
// The enqueuer's NewMessage() method is called internally to create the message object.
// The created message object implements the api.Message[R] interface.
func (q *queue[R]) NewMessage() api.Message[R] {
	return q.enqueuer.NewMessage()
}

// Enqueue adds an item of type M to the queue, following the context.
// It delegates the operation to its enqueuer and returns any error produced during this operation.
func (q *queue[R]) Enqueue(ctx context.Context, msg api.Message[R]) error {
	return q.enqueuer.Enqueue(ctx, msg)
}

// Dequeue retrieves and removes an item from the queue, following the context.
// It delegates the operation to its dequeuer and returns the dequeued item along with any error that occurred during the operation.
func (q *queue[R]) Dequeue(ctx context.Context) (api.DequeueMessage[R], error) {
	return q.dequeuer.Dequeue(ctx)
}

// Disconnect disconnects from the queue by calling the `Disconnect()` method on both the enqueuer and the dequeuer.
// It delegates the disconnection operations to both interfaces concurrently and waits for them to complete.
// It returns any error occurred during the disconnection.
func (q *queue[R]) Disconnect(ctx context.Context) error {

	errGrp := new(errgroup.Group)

	errGrp.Go(func() error {
		return q.enqueuer.Disconnect(ctx)
	})
	errGrp.Go(func() error {
		return q.dequeuer.Disconnect(ctx)
	})

	return errGrp.Wait()
}
