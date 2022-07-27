package ticker

import (
	"time"
)

type Ticker string

const (
	BTCUSDTicker Ticker = "BTC_USD"
)

type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string // decimal value. example: "0", "10", "12.2", "13.2345122"
}

type PriceStreamSubscriber interface {
	SubscribePriceStream(Ticker) (chan TickerPrice, chan error)
	// Added ProviderID() method to the subscriber interface.
	// It might be useful to distinguish different subscribers while calculating average ticker price.
	ProviderID() int
	// Added Stop() method to properly close subscription and notify data provider that application stopped reads stream.
	Stop()
}
