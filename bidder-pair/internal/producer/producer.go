package producer

import (
	"context"
	"log"
	"sync"
	"time"
)

type ProducerConfig struct {
	Logger        *log.Logger
	FlushInterval time.Duration
	SlowSend      time.Duration
}

type Producer struct {
	logger *log.Logger

	mu sync.Mutex

	flushInterval time.Duration
	slowSend      time.Duration
}

func NewProducer(cfg ProducerConfig) *Producer {
	p := &Producer{
		logger:        cfg.Logger,
		flushInterval: cfg.FlushInterval,
		slowSend:      cfg.SlowSend,
	}
	return p
}

func (p *Producer) Send(e Event) {
	go func() {
		// slow send
		time.Sleep(p.slowSend)
	}()
}

func (p *Producer) Close(ctx context.Context) error {
	return nil
}
