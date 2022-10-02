package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lfcamarati/pfa-go/internal/order/infra/database"
	"github.com/lfcamarati/pfa-go/internal/order/usecase"
	"github.com/lfcamarati/pfa-go/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/orders")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	orderRepository := database.NewOrderRepository(db)
	calculateFinalPriceUseCase := usecase.NewCalculateFinalPriceUseCase(orderRepository)
	ch, err := rabbitmq.OpenChannel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// Wait Groups to control how many works will read the messages from rabbitmq
	maxWorkers := 5
	wg := sync.WaitGroup{}

	out := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, out)

	wg.Add(maxWorkers)
	for i := 1; i <= maxWorkers; i++ {
		defer wg.Done()
		go worker(out, calculateFinalPriceUseCase, i)
	}
	wg.Wait()

	// orderInput := usecase.OrderInputDTO{
	// 	ID:    "456",
	// 	Price: 100.0,
	// 	Tax:   10.0,
	// }

	// orderOutput, err := calculateFinalPriceUseCase.Execute(orderInput)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("The final price is %f", orderOutput.FinalPrice)
}

func worker(deliveryMessage <-chan amqp.Delivery, uc *usecase.CalculateFinalPriceUseCase, workderId int) {
	for msg := range deliveryMessage {
		var orderInput usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &orderInput)

		if err != nil {
			fmt.Println("Error unmarshalling message", err)
		}

		orderInput.Tax = 10.0
		_, err = uc.Execute(orderInput)

		if err != nil {
			fmt.Println("Error executing usecase", err)
		}

		msg.Ack(false)
		fmt.Println("Worker", workderId, "processed order", orderInput.ID)
	}
}
