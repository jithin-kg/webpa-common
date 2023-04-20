package mhook

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// Service describes the core operations around webhook subscriptions.
type Service interface {
	// Add adds the given owned webhook to the current list of webhooks. If the operation
	// succeeds, a non-nil error is returned.
	Add(owner string, w *Webhook) error

	// AllWebhooks lists all the current webhooks for the given owner.
	// If an owner is not provided, all webhooks are returned.
	AllWebhooks(owner string) ([]*Webhook, error)
}
type loggerGroup struct {
	Error log.Logger
	Debug log.Logger
}
type service struct {
	store WebhookStore
	// loggers *loggerGroup
}

func (s *service) Add(owner string, w *Webhook) error {
	fmt.Printf("__Service.go: Add() called with owner %s and webhook %+v\n", owner, w)
	return s.store.Add(owner, w)
}

func (s *service) AllWebhooks(owner string) ([]*Webhook, error) {
	// s.loggers.Debug.Log("msg", "AllWebhooks called", "owner", owner)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic:", r)
			debug.PrintStack()
		}
	}()

	webhooks, err := s.store.AllWebhooks(owner)
	fmt.Println("__")
	if err != nil {
		return nil, err
	}

	// Create a new slice of *Webhook type
	webhooksPtr := make([]*Webhook, len(webhooks))

	// Copy the contents of webhooks slice to webhooksPtr slice
	for i, wh := range webhooks {
		webhooksPtr[i] = wh
	}

	return webhooksPtr, nil
}

func Initialize(cfg *Config, watches ...Watch) (Service, func(), error) {

	store := NewWebhookStore()

	svc := &service{
		// loggers: newLoggerGroup(logger),
		store: store,
	}

	// ...

	return svc, func() { /*...*/ }, nil
}

func newLoggerGroup(root log.Logger) *loggerGroup {
	if root == nil {
		root = log.NewNopLogger()
	}

	return &loggerGroup{
		Debug: log.WithPrefix(root, level.Key(), level.DebugValue()),
		Error: log.WithPrefix(root, level.Key(), level.ErrorValue()),
	}

}
func validateConfig(cfg *Config) {
	if cfg.WatchUpdateInterval == 0 {
		cfg.WatchUpdateInterval = time.Second * 5
	}
}
