// Package api defines several abstract interfaces for enqueue and dequeue operations in a queuing system.
// These interfaces are designed to be implemented by specific internal packages for different systems (like Oracle Advanced Queues)
// and are exposed to the user via applicable packages under ezQue. The specific types are yet to be determined (the use of Go generics).
//
// Dequeuer interface defines the Dequeue() and Disconnect() methods. Dequeue() blocks until a message is read from the queue,
// while Disconnect() is used to cut the connection with the queue.
//
// DequeueMessage interface represents a message that has been dequeued. It provides methods to retrieve the message itself, acknowledge it,
// and reject its acknowledgment.
//
// Enqueuer interface describes the Enqueue() and Disconnect() methods. Enqueue() pushes a new message onto the queue,
// and Disconnect() ends the connection with the queue.
//
// Message interface is a representation of a general message that can be enqueued or dequeued. It provides methods
// for accessing the raw message and its text representation.
//
// The interfaces defined in this package serve as a contract for any system-specific implementation ensuring interoperability
// and consistent usage.
package api
