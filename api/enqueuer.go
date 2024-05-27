package api

import "context"

// Enqueuer is an interface that provides methods for managing enqueue operations to an Oracle Advanced Queue.
// It defines the NewMessage() method for creating a new message object, the Enqueue() method for
// enqueueing a message to the queue, and the Disconnect() method for disconnecting from the queue.
type Enqueuer[R any] interface {
	NewMessage() Message[R]
	Enqueue(ctx context.Context, msg Message[R]) error
	Disconnect(ctx context.Context) error
}
