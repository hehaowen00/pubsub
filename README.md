# pubsub

Golang implementation of publish-subscribe pattern using channels.

```go
err := pubsub.NewTopic[int]("topic", pubsub.NewRingBuffer[int](16))
if err != nil {
    // ...
}

subscriber, err = pubsub.Subscribe[int]("topic")
if err != nil {
    // ...
}

publisher, err = pubsub.NewPublisher[int]("topic")
if err != nil {
    // ...
}
```

- Topics are explicitly created instead of implicit

- Messages are not persisted to disk

- Each subscriber has its own internal buffer to prevent deadlocks

- Does not deep copy structs
