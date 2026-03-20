package producer

import (
	"context"
	"log"
	"time"
)

type ProducerConfig struct {
	Logger        *log.Logger
	FlushInterval time.Duration
	SlowSend      time.Duration
}

type Producer struct {
	logger *log.Logger

	flushInterval time.Duration
}

func NewProducer(cfg ProducerConfig) *Producer {
	p := &Producer{
		logger:        cfg.Logger,
		flushInterval: cfg.FlushInterval,
	}
	return p
}

func (p *Producer) Enqueue(event *Event) {
}

func (p *Producer) SendToGCP(event []*Event) {
	go func() {
		// slow send
		time.Sleep(500 * time.Millisecond)
	}()
}

func (p *Producer) Close(ctx context.Context) error {
	return nil
}
