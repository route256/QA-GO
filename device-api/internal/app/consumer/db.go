package consumer

import (
	"context"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/repo"
	"sync"
	"time"

	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/model"
)

// Consumer is a consumer of data from repository
type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.DeviceEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	done chan bool
	wg   *sync.WaitGroup
}

// NewDbConsumer creates new DB consumer
func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.DeviceEvent) Consumer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
		done:      done,
	}
}

func (c *consumer) Start() {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(context.TODO(), c.batchSize)
					if err != nil {
						continue
					}
					for _, event := range events {
						c.events <- event
					}
				case <-c.done:
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	close(c.done)
	c.wg.Wait()
}
