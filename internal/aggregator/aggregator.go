// Package aggregator contains simple object that calculates single ticker price from given array.
// I've decided not try to implement any complex aggregating fuction.
// Instead of it this implementation calculates average price per provider
// and then takes a median value of it.
package aggregator

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/adanekin/occamfi-task/internal/ticker"
)

type ProviderMinutePrice struct {
	ProviderID int
	Prices     []ticker.TickerPrice
}

type Aggregator struct{}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) Aggregate(prices []ProviderMinutePrice) float64 {
	avgProviderPrices := make([]float64, 0, len(prices))
	for _, price := range prices {
		avgProviderPrice, err := avgProviderPrice(price.Prices)
		if err != nil {
			fmt.Printf("failed to calculate price: %v\n", err)
			continue
		}

		avgProviderPrices = append(avgProviderPrices, avgProviderPrice)
	}

	return median(avgProviderPrices)
}

func avgProviderPrice(prices []ticker.TickerPrice) (float64, error) {
	var sum float64
	for _, price := range prices {
		val, err := strconv.ParseFloat(price.Price, 64)
		if err != nil {
			return 0, err
		}

		sum += val
	}

	return sum / float64(len(prices)), nil
}

func median(arr []float64) float64 {
	sort.Float64s(arr)

	l := len(arr)
	if l%2 == 0 {
		return (arr[l/2-1] + arr[l/2]) / 2
	}

	return arr[l/2]
}
