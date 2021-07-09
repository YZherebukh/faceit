//go:generate mockgen -source ../notifier/notifier.go -destination ../notifier/mock/mock_notifier.go

package notifier

import (
	"context"
	"encoding/json"

	"github.com/faceit/test/logger"
)

type notifier interface {
	Send(consumer []string, message []byte) error
}

// Notifier is a notifier struct
type Notifier struct {
	notifier
	consumers    []string
	consumersMap map[string]struct{}
	log          logger.Logger
}

func New(n notifier, consumers []string, l logger.Logger) *Notifier {
	notifier := &Notifier{
		notifier: n,
		log:      l,
	}

	notifier.consumersMap = make(map[string]struct{})

	for _, c := range consumers {
		if _, ok := notifier.consumersMap[c]; !ok {
			notifier.consumersMap[c] = struct{}{}
			notifier.consumers = append(notifier.consumers, c)
		}
	}

	return notifier
}

// Do sends a messages to one or many consumers
// if consumers slice is empty, would send messages to all consumers
func (n *Notifier) Do(ctx context.Context, message interface{}, consumers []string) {
	messageByte, err := json.Marshal(message)
	if err != nil {
		n.log.Errorf(ctx, "failed to marshal message %#v, error: %w", message, err)
		return
	}

	cc := n.getConsumers(consumers)
	if len(cc) == 0 {
		n.log.Warningf(ctx, "no consumers to send message to")
		return
	}

	err = n.Send(cc, messageByte)
	if err != nil {
		n.log.Errorf(ctx, "failed to send message to %#v, error: %w", cc, err)
	}
}

func (n *Notifier) getConsumers(cc []string) []string {
	if len(cc) < 0 {
		return n.consumers
	}

	consumers := make([]string, 0, len(n.consumers))

	for _, c := range cc {
		if _, ok := n.consumersMap[c]; ok {
			consumers = append(consumers, c)
		}
	}

	return consumers
}
