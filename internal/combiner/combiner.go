// Package coombiner provides main logic for combining ticker prices from multiple streams.
package combiner

import (
	"fmt"
	"sync"
	"time"

	"github.com/adanekin/occamfi-task/internal/aggregator"
	"github.com/adanekin/occamfi-task/internal/ticker"
)

// Interface for object that calculates single ticker price from all providers by given price arrays.
type Aggregator interface {
	Aggregate(prices []aggregator.ProviderMinutePrice) float64
}

// Interface for sending calculated prices.
type Sender interface {
	Send(timestamp time.Time, price float64) error
}

type Combiner struct {
	current        map[int][]ticker.TickerPrice
	pending        map[time.Time]map[int][]ticker.TickerPrice
	providersCount int
	aggregator     Aggregator
	sender         Sender
	mu             sync.Mutex
}

func NewCombiner(providerIDs []int, aggregator Aggregator, sender Sender) *Combiner {
	return &Combiner{
		current:        make(map[int][]ticker.TickerPrice),
		pending:        make(map[time.Time]map[int][]ticker.TickerPrice),
		providersCount: len(providerIDs),
		aggregator:     aggregator,
		sender:         sender,
		mu:             sync.Mutex{},
	}
}

func (c *Combiner) AddPrice(price ticker.TickerPrice, providerID int) {
	min := minute(price.Time)

	c.mu.Lock()

	cur, ok := c.current[providerID]
	if !ok {
		c.current[providerID] = []ticker.TickerPrice{price}

		return
	}

	if cur[len(cur)-1].Time.Minute() == price.Time.Minute() {
		c.current[providerID] = append(c.current[providerID], price)

		return
	}

	if _, ok := c.pending[min]; !ok {
		c.pending[min] = make(map[int][]ticker.TickerPrice)
	}

	c.pending[min][providerID] = cur
	c.current[providerID] = []ticker.TickerPrice{price}

	c.mu.Unlock()

	if len(c.pending[min]) == c.providersCount {
		aggPrices := make([]aggregator.ProviderMinutePrice, 0)
		for provID, prices := range c.pending[min] {
			aggPrices = append(aggPrices, aggregator.ProviderMinutePrice{ProviderID: provID, Prices: prices})
		}

		if err := c.sender.Send(min, c.aggregator.Aggregate(aggPrices)); err != nil {
			fmt.Printf("failed to send price: %v", err)
			return
		}
	}
}

func minute(ts time.Time) time.Time {
	return time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), 0, 0, time.UTC)
}
