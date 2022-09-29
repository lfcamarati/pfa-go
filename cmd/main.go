package main

import (
	"database/sql"
	"fmt"

	"github.com/lfcamarati/pfa-go/internal/order/infra/database"
	"github.com/lfcamarati/pfa-go/internal/order/usecase"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(mysql:3306)/orders")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	orderRepository := database.NewOrderRepository(db)
	calculateFinalPriceUseCase := usecase.NewCalculateFinalPriceUseCase(orderRepository)

	orderInput := usecase.OrderInputDTO{
		ID:    "456",
		Price: 100.0,
		Tax:   10.0,
	}

	orderOutput, err := calculateFinalPriceUseCase.Execute(orderInput)

	if err != nil {
		panic(err)
	}

	fmt.Printf("The final price is %f", orderOutput.FinalPrice)
}
