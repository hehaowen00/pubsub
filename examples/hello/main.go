package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hehaowen00/pubsub"
)

func subscribeLoop() {
	sub, err := pubsub.NewSubscriber[int]("topic")
	if err != nil {
		log.Println("error unable to subscribe to topic: ", err)
		return
	}

	for {
		_, ok := <-sub.Recv()
		if !ok {
			fmt.Println("closed")
			break
		}

		data := sub.Read()

		fmt.Println("recv", data)
	}
}

func main() {
	err := pubsub.NewTopic[int]("topic")
	if err != nil {
		panic(err)
	}

	go subscribeLoop()

	time.Sleep(time.Second)

	publisher, err := pubsub.NewPublisher[int]("topic")
	if err != nil {
		panic(err)
	}

	log.Println(publisher.Publish(1))
	log.Println(publisher.Publish(2))

	go subscribeLoop()

	time.Sleep(time.Second)

	log.Println(publisher.Publish(3))
	log.Println(publisher.Publish(4))

	pubsub.Close[int]("topic")

	log.Println(publisher.Publish(5))
	log.Println(publisher.Publish(6))

	fmt.Println("Waiting for interrupt...")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
