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
