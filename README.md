# ezQue

ezQue is a simple Go interface that abstracts various messaging queue systems like Oracle Advance Queues (OracleAQ), Apache Kafka, Apache ActiveMQ, and others. ezQue aims to provide simple queue access APIs, making it easier to interact with different messaging systems in a unified way.

Even though ezQue's long-term goal is to support many queue systems, its current release only supports OracleAQ. Upcoming releases plan to include support for ActiveMQ/Artemis, followed by Apache Kafka.

With ezQue, you can easily Enqueue/Dequeue messages from the queue. One important point to note is that after Dequeue, messages should be acknowledged (Ack) or not acknowledged (Ack) to properly manage the message lifecycle.

## Installation

You can get ezQue by using:

```sh
go get github.com/pgvanniekerk/ezQue
```

Ensure your Go project is using Go Modules (it will have a `go.mod` file in its root) and that your `GO111MODULE` environment variable is either `auto`, or `on`.

## Connecting to Oracle

You can establish a connection to Oracle Advance Queues (OracleAQ) via ezQue's `Connect` function. The `Connect` function requires a `queueConnector` function and certain connection parameters as `Options`, which include authentication and database details. Here's an example on how to use the `Connect` function:

```go
package main

import (
    "context"
    "time"

    "github.com/pgvanniekerk/ezQue"
)

func main() {
    ctx := context.Background()

    server := os.Getenv("DB_SERVER")
    portString := os.Getenv("DB_PORT")
    port, err := strconv.Atoi(portString)
    if err != nil {
        t.Fatalf("Invalid port number: %v", err)
    }
    username := os.Getenv("DB_USERNAME")
    password := os.Getenv("DB_PASSWORD")
    sid := os.Getenv("DB_SID")
    
    q, err := ezQue.Connect(oraaq.OracleAqJms,
        oraaq.Queue("text_msg_queue",
            oraaq.LocatedAt(server, uint16(port)),
            oraaq.AuthenticatedWith(username, password),
            oraaq.UsingSID(sid),
        ),
    )
    if err != nil {
        t.Error(err)
    }

    // Continued code block.
	
    err = q.Disconnect(context.Background())
    if err != nil {
        t.Error(err)
    }

}
```

## Dequeueing Messages

You can dequeue messages from the established connection using the Dequeue method. Below is a simple example demonstrating how to dequeue messages:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time

	"github.com/pgvanniekerk/ezQue"
)

func main() {
    // create the context
    ctx := context.Background()
    
    /* initialize your q here */
    
    // Dequeue messages
    for {
        // Dequeue a message from the queue
        dequeueMessage, err := q.Dequeue(ctx)
        if err != nil {
            log.Fatalf("Dequeueing failed: %v", err)
        }
    
        // Get the raw message
        // The Raw() method provides access to the underlying queue message implementation.
        rawMessage := dequeueMessage.Message().Raw()
    
        // Process your raw message...
        fmt.Printf("Received raw: %v\n", rawMessage)
    
        // Get the text message
        // Use the Text() method to extract the critical text information from the message
        textMessage := dequeueMessage.Message().Text()
    
        // Process your text message...
        fmt.Printf("Received text: %v\n", textMessage)
    
        // Acknowledge the message.
        err = dequeueMessage.Ack(ctx)
        if err != nil {
            log.Fatalf("Failed to acknowledge: %v", err)
        }
    
        // Take a break for the next dequeue operation.
        time.Sleep(time.Second)
    }
    
    // Don't forget to disconnect after finished
    err := q.Disconnect(ctx)
    if err != nil {
        log.Fatalf("disconnection failed: %v", err)
    }
}
```

## Enqueueing Messages

You can enqueue messages to the established connection using the Enqueue method. Below is a simple example demonstrating how to enqueue messages:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pgvanniekerk/ezQue"
	"github.com/pgvanniekerk/ezQue/api"
)

func main() {
    // create the context
    ctx := context.Background()
    
    /* initialize your q here */
    
    // Enqueue messages
    for i := 0; i < 10; i++ {
        // Create your message
        msg := q.NewMessage()
        
        // Enqueue a message to the queue
        err := q.Enqueue(ctx, msg)
        if err != nil {
            log.Fatalf("Enqueueing failed: %v", err)
        }
        fmt.Printf("Enqueued: %v\n", msg)
    }
    
    // Do not forget to disconnect
    err := q.Disconnect(ctx)
    if err != nil {
        log.Fatalf("disconnection failed: %v", err)
    }
}
```