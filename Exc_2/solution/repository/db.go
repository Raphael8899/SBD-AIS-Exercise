package repository

import (
	"ordersystem/model"
	"time"
)

type DatabaseHandler struct {
	// drinks represent all available drinks
	drinks []model.Drink
	// orders serves as order history
	orders []model.Order
}

func NewDatabaseHandler() *DatabaseHandler {
	// Init the drinks slice with some test data
	drinks := []model.Drink{
		{
			ID:          1,
			Name:        "Beer",
			Price:       2.00,
			Description: "Hagenberger Gold",
		},
		{
			ID:          2,
			Name:        "Spritzer",
			Price:       1.40,
			Description: "Wine with soda",
		},
		{
			ID:          3,
			Name:        "Coffee",
			Price:       0.00,
			Description: "Mifare isn't that secure ;)",
		},
	}

	// Init orders slice with some test data
	orders := []model.Order{
		{
			DrinkID:   1,
			Amount:    2,
			CreatedAt: time.Now().Add(-time.Hour * 2),
		},
		{
			DrinkID:   2,
			Amount:    1,
			CreatedAt: time.Now().Add(-time.Hour),
		},
	}

	return &DatabaseHandler{
		drinks: drinks,
		orders: orders,
	}
}

func (db *DatabaseHandler) GetDrinks() []model.Drink {
	return db.drinks
}

func (db *DatabaseHandler) GetOrders() []model.Order {
	return db.orders
}

func (db *DatabaseHandler) GetTotalledOrders() map[uint64]uint64 {
	// calculate total orders
	// key = DrinkID, value = Amount of orders
	totalledOrders := make(map[uint64]uint64)
	
	// loop through all orders and sum up the amounts per drink
	for _, order := range db.orders {
		totalledOrders[order.DrinkID] += order.Amount
	}
	
	return totalledOrders
}

func (db *DatabaseHandler) AddOrder(order *model.Order) {
	// add order to db.orders slice
	db.orders = append(db.orders, *order)
}

