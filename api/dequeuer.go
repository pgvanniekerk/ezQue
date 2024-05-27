package api

import (
	"context"
)

// Dequeuer provides the Dequeue() method, popping a message of type M and wrapping
// it as a DequeueMessage. Calling Dequeue() will block until a message has been read
// from the bound queue.
type Dequeuer[R any] interface {
	Dequeue(ctx context.Context) (DequeueMessage[R], error)
	Disconnect(ctx context.Context) error
}

// DequeueMessage represents a message in the deque. It provides methods to
// retrieve the message (Message()), acknowledge it (Ack()), and negate
// its acknowledgment (NAck()). Types representing deque messages should
// implement this interface.
type DequeueMessage[R any] interface {
	Message() Message[R]
	Ack(ctx context.Context) error
	NAck(ctx context.Context) error
}
