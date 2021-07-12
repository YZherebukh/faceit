//go:generate mockgen -source ../queue/queue.go -destination ../queue/mock/mock_queue.go

package queue

import (
	"context"
	"time"

	"github.com/faceit/test/config"
	"github.com/faceit/test/entity"
)

type notifier interface {
	Do(ctx context.Context, consumers []string, message interface{})
}

// Queue is a queue implementation
type Queue struct {
	addCh      chan struct{}
	popCh      chan struct{}
	addMessage chan entity.NotifierMessage
	popMessage chan entity.NotifierMessage
	notifier   notifier
}

// New creates new queue
func New(cfg config.Queue, n notifier) *Queue {
	q := &Queue{
		addCh:      make(chan struct{}, cfg.QueueSize),
		popCh:      make(chan struct{}, cfg.GoRoutinesSize),
		addMessage: make(chan entity.NotifierMessage, 1),
		popMessage: make(chan entity.NotifierMessage, 1),
		notifier:   n,
	}

	go q.pop()
	go q.add()

	return q
}

// Add adds new notifyReq to queue
func (q Queue) Add(message entity.NotifierMessage) {
	q.addCh <- struct{}{}

	go func(message entity.NotifierMessage) {
		q.addMessage <- message
	}(message)
}

func (q Queue) add() {
	for message := range q.addMessage {
		<-q.addCh

		q.popMessage <- message
	}
}

// pop pops item from queue
func (q Queue) pop() {
	for message := range q.popMessage {
		q.popCh <- struct{}{}

		go func(message entity.NotifierMessage) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Minute))
			defer cancel()
			q.notifier.Do(ctx, message.Consumers, message.Message)
			<-q.popCh
		}(message)
	}
}

// Close closing queue add channel
func (q Queue) Closed() {
	for len(q.addCh) != 0 && len(q.popCh) != 0 {

	}

	close(q.addMessage)
	close(q.popMessage)
}
