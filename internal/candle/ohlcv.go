package candle

import "time"

type Candle struct {
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

type OHLCVAggregator struct {
	OneSecCandles []Candle
	OneMinCandles []Candle
	last1sTime    time.Time
	last1mTime    time.Time
}

func NewOHLCVAggregator() *OHLCVAggregator {
	return &OHLCVAggregator{}
}

func (o *OHLCVAggregator) AddTrade(price, quantity float64, t time.Time) {
	o.update1s(price, quantity, t)
	o.update1m(price, quantity, t)
}

func (o *OHLCVAggregator) update1s(price, qty float64, t time.Time) {
	ts := t.Truncate(time.Second)
	if ts != o.last1sTime {
		o.last1sTime = ts
		o.OneSecCandles = append(o.OneSecCandles, Candle{
			Open:      price,
			High:      price,
			Low:       price,
			Close:     price,
			Volume:    qty,
			Timestamp: ts,
		})
	} else {
		c := &o.OneSecCandles[len(o.OneSecCandles)-1]
		if price > c.High {
			c.High = price
		}
		if price < c.Low {
			c.Low = price
		}
		c.Close = price
		c.Volume += qty
	}
}

func (o *OHLCVAggregator) update1m(price, qty float64, t time.Time) {
	ts := t.Truncate(time.Minute)
	if ts != o.last1mTime {
		o.last1mTime = ts
		o.OneMinCandles = append(o.OneMinCandles, Candle{
			Open:      price,
			High:      price,
			Low:       price,
			Close:     price,
			Volume:    qty,
			Timestamp: ts,
		})
	} else {
		c := &o.OneMinCandles[len(o.OneMinCandles)-1]
		if price > c.High {
			c.High = price
		}
		if price < c.Low {
			c.Low = price
		}
		c.Close = price
		c.Volume += qty
	}
}
