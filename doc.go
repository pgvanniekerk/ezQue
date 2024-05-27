// Package ezQue proposes a simplified Go interface for working with various messaging queue systems.
// The package is designed to abstract the complexities of different queue systems such as Oracle Advance Queues, Apache Kafka, Apache ActiveMQ
// and others, providing a unified experience across all supported systems.
//
// Specifically, ezQue provides the Queue interface with essential queue operations — Enqueue, Dequeue, and Disconnect.
// The Enqueue method pushes a new message onto the queue, the Dequeue method retrieves a message from the queue
// and returns it wrapped with a DequeueMessage (this ensures message processing can be acknowledged or not acknowledged appropriately),
// and the Disconnect method terminates the connection with the queue.
//
// A queueConnector is provided as a function type for connecting to a specific queue system, configuring the connection options,
// and initializing Enqueuer and Dequeuer interfaces.
//
// The Connect function initializes a connection to a queue, taking as parameters a queueConnector for getting
// system-specific Enqueuer and Dequeuer and an options parameter for Enqueuer/Dequeuer configuration.
//
// Note: The current release of ezQue only supports OracleAQ but the design intends to accommodate additional queue systems
// such as ActiveMQ/Artemis and Apache Kafka in the future.
//
// By making queuing operations consistent and system-agnostic, ezQue makes it easier for developers to implement
// various operations on any supported messaging system.
package ezQue
