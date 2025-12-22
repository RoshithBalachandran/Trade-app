package engine

import "time"

type Trade struct {
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	Timestamp   time.Time `json:"timestamp"`
}

func (ob *OrderBook) Match(order *Order) []*Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	trades := []*Trade{}

	if order.Side == Buy {
		for i := 0; i < len(ob.Sells) && order.Quantity > 0; i++ {
			sell := ob.Sells[i]
			if order.Type == Limit && order.Price < sell.Price {
				break
			}
			qty := min(order.Quantity, sell.Quantity)
			trade := &Trade{
				BuyOrderID:  order.ID,
				SellOrderID: sell.ID,
				Price:       sell.Price,
				Quantity:    qty,
				Timestamp:   time.Now(),
			}
			trades = append(trades, trade)
			order.Quantity -= qty
			sell.Quantity -= qty
			if sell.Quantity == 0 {
				ob.Sells = append(ob.Sells[:i], ob.Sells[i+1:]...)
				i--
			}
		}
		if order.Quantity > 0 && order.Type == Limit {
			ob.Buys = append(ob.Buys, order)
		}
	} else {
		for i := 0; i < len(ob.Buys) && order.Quantity > 0; i++ {
			buy := ob.Buys[i]
			if order.Type == Limit && order.Price > buy.Price {
				break
			}
			qty := min(order.Quantity, buy.Quantity)
			trade := &Trade{
				BuyOrderID:  buy.ID,
				SellOrderID: order.ID,
				Price:       buy.Price,
				Quantity:    qty,
				Timestamp:   time.Now(),
			}
			trades = append(trades, trade)
			order.Quantity -= qty
			buy.Quantity -= qty
			if buy.Quantity == 0 {
				ob.Buys = append(ob.Buys[:i], ob.Buys[i+1:]...)
				i--
			}
		}
		if order.Quantity > 0 && order.Type == Limit {
			ob.Sells = append(ob.Sells, order)
		}
	}

	return trades
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
