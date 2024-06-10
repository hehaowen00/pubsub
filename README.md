# pubsub

Golang implementation of publish-subscribe pattern using channels.

```go
sub, err = pubsub.Subscribe[int]("topic")
if err != nil {
    // ...
}


err = pubsub.Publish("topic", 1)
if err != nil {
    // ...
}
```

## Details

- The message broker keeps the last 16 messages in memory

- Messages are not persisted to disk

- Each subscriber receive channel is buffered

- Subscribers are expected to consume messages or the broker
will deadlock trying to send to a full channel
