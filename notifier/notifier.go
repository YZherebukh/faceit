//go:generate mockgen -source ../notifier/notifier.go -destination ../notifier/mock/mock_notifier.go

package notifier

import (
	"context"
	"encoding/json"
	"time"

	"github.com/faceit/test/config"
	"github.com/faceit/test/logger"
)

type notifier interface {
	Send(ctx context.Context, consumer []string, message []byte) error
}

// Notifier is a notifier struct
type Notifier struct {
	notifier              notifier
	timeout               int
	clientMaxRetry        int
	clientTimeoutIncrease int
	consumers             []string
	consumersMap          map[string]struct{}
	log                   logger.Logger
}

func New(cfg config.Notifier, n notifier, consumers []string, l logger.Logger) *Notifier {
	notifier := &Notifier{
		notifier:              n,
		log:                   l,
		timeout:               cfg.Timeout,
		clientMaxRetry:        cfg.ClientMaxRetry,
		clientTimeoutIncrease: cfg.ClientTimeoutIncrease,
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
func (n *Notifier) Do(ctx context.Context, consumers []string, message interface{}) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(n.timeout))
	defer cancel()

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

	err = n.sendWithRetry(ctx, cc, messageByte)
	if err != nil {
		n.log.Errorf(ctx, "after %d retries still failed to send message, error: ", err)
	}
}

func (n *Notifier) getConsumers(cc []string) []string {
	if len(cc) == 0 {
		return cc
	}

	consumers := make([]string, 0, len(n.consumers))

	for _, c := range cc {
		if _, ok := n.consumersMap[c]; ok {
			consumers = append(consumers, c)
		}
	}

	return consumers
}

func (n *Notifier) sendWithRetry(ctx context.Context, consumers []string, message []byte) error {
	var (
		i   int
		err error
	)

	for i < n.clientMaxRetry {
		err = n.notifier.Send(ctx, consumers, message)
		if err != nil {
			n.log.Warningf(ctx, "failed to end message to %#v, error: %w", consumers, err)

			time.Sleep(time.Duration(n.clientTimeoutIncrease*i) * time.Second)

			i++
			continue
		}

		return nil
	}

	return err
}
