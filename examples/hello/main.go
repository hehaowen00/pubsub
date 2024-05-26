package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hehaowen00/pubsub"
)

func main() {
	go func() {
		sub, err := pubsub.Subscribe[int]("topic")
		if err != nil {
			log.Println("error unable to subscribe", err)
			return
		}

		for {
			m, ok := <-sub.Recv()
			if !ok {
				fmt.Println("closed")
				break
			}

			fmt.Println("recv", m)
		}
	}()

	time.Sleep(time.Second)

	try(pubsub.Publish("topic", 1))
	try(pubsub.Publish("topic", 2))

	go func() {
		sub, err := pubsub.Subscribe[int]("topic")
		if err != nil {
			log.Println("error unable to subscribe", err)
			return
		}

		for {
			m, ok := <-sub.Recv()
			if !ok {
				fmt.Println("closed")
				break
			}

			fmt.Println("recv", m)
		}
	}()

	time.Sleep(time.Second)

	try(pubsub.Publish("topic", 3))
	try(pubsub.Publish("topic", 4))

	pubsub.Close[int]("topic")

	try(pubsub.Publish("topic", 5))
	try(pubsub.Publish("topic", 6))

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}

func try(err error) {
	if err != nil {
		panic(err)
	}
}
