package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Order struct {
	ID    string
	Price float64
}

func generateOrders() Order {
	return Order{
		ID:    uuid.New().String(),
		Price: rand.Float64() * 100.0,
	}
}

func Notify(ch *amqp.Channel, order Order) error {
	body, err := json.Marshal(order)

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		context.TODO(),
		"amq.direct",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

func main() {
	rand.Seed(time.Now().UnixNano())

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	for i := 0; i < 100000; i++ {
		order := generateOrders()
		err := Notify(ch, order)
		if err != nil {
			panic(err)
		}
		// fmt.Println(order)
	}
}
