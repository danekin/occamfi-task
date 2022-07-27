// Package receiver takes responsibility for subscribing on price streams,
// resubscribing in case of errors and sending data from streams to Combiner
package receiver

import (
	"context"
	"sync"

	"github.com/adanekin/occamfi-task/internal/ticker"
)

type Combiner interface {
	AddPrice(price ticker.TickerPrice, subscriberID int)
}

type PriceReceiver struct {
	subs     []ticker.PriceStreamSubscriber
	combiner Combiner
}

func NewPriceReceiver(combiner Combiner, subs []ticker.PriceStreamSubscriber) *PriceReceiver {
	return &PriceReceiver{
		combiner: combiner,
		subs:     subs,
	}
}

func (r *PriceReceiver) Run(ctx context.Context) {
	wg := sync.WaitGroup{}

	for _, sub := range r.subs {
		wg.Add(1)
		go func(subscriber ticker.PriceStreamSubscriber) {
			r.subscribe(ctx, subscriber)
			wg.Done()
		}(sub)
	}

	wg.Wait()
}

func (r *PriceReceiver) subscribe(ctx context.Context, sub ticker.PriceStreamSubscriber) {
	dataChan, errChan := sub.SubscribePriceStream(ticker.BTCUSDTicker)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		for price := range dataChan {
			r.combiner.AddPrice(price, sub.ProviderID())
		}
		wg.Done()
	}()

	select {
	case <-errChan:
		r.subscribe(ctx, sub)
	case <-ctx.Done():
		// Stops combiner goroutine.
		sub.Stop()
		wg.Wait()
		return
	}
}
