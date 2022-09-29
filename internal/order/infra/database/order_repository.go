package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lfcamarati/pfa-go/internal/order/entity"
)

type OrderDatabaseRepository struct {
	Db *sql.DB
}

// CREATE TABLE `orders` (
// 	`id` varchar(255) NOT NULL,
// 	`price` float NOT NULL,
// 	`tax` float NOT NULL,
// 	`final_price` float NOT NULL,
// 	PRIMARY KEY (`id`))

func NewOrderRepository(db *sql.DB) *OrderDatabaseRepository {
	return &OrderDatabaseRepository{Db: db}
}

func (orderRepository *OrderDatabaseRepository) Save(order *entity.Order) error {
	stmt, err := orderRepository.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")

	if err != nil {
		return err
	}

	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)

	if err != nil {
		return err
	}

	return nil
}
