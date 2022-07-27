package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/adanekin/occamfi-task/internal/aggregator"
	"github.com/adanekin/occamfi-task/internal/combiner"
	"github.com/adanekin/occamfi-task/internal/receiver"
	"github.com/adanekin/occamfi-task/internal/sender"
	"github.com/adanekin/occamfi-task/internal/ticker"
)

func main() {
	file, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open file: %v\n", err)
		return
	}

	subs, providerIDs := generateMockSubscribers()

	writer := bufio.NewWriter(file)
	fileSender := sender.NewFileSender(writer)
	priceCombiner := combiner.NewCombiner(
		providerIDs,
		aggregator.NewAggregator(),
		fileSender,
	)

	priceReceiver := receiver.NewPriceReceiver(
		priceCombiner,
		subs,
	)

	stopChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		priceReceiver.Run(ctx)
		wg.Done()
	}()

	fmt.Printf("application started\n")

	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	<-stopChan

	cancel()
	wg.Wait()

	file.Sync()
	file.Close()

	fmt.Printf("application stopped\n")
}

type MockSubscriber struct {
	providerID int
}

func NewMockSubscriber(providerID int) *MockSubscriber {
	return &MockSubscriber{
		providerID: providerID,
	}
}

func (s *MockSubscriber) SubscribePriceStream(ticker.Ticker) (chan ticker.TickerPrice, chan error) {
	return nil, nil
}

func (s *MockSubscriber) ProviderID() int {
	return s.providerID
}

func (s *MockSubscriber) Stop() {}

func generateMockSubscribers() ([]ticker.PriceStreamSubscriber, []int) {
	ids := make([]int, 0, 100)
	subs := make([]ticker.PriceStreamSubscriber, 0, 100)

	for i := 0; i < 100; i++ {
		ids = append(ids, i)
		subs = append(subs, NewMockSubscriber(i))
	}

	return subs, ids
}
