package engine

import "sync"


type OrderBook struct{
	Buys []*Order
	Sells []*Order
	mu sync.Mutex
}

func NewOrderBook()*OrderBook{
	return &OrderBook{
		Buys: []*Order{},
		Sells: []*Order{},
	}
}