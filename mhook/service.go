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
	Add(owner string, w *Webhook) error
	AllWebhooks(owner string) ([]*Webhook, error)
}

type loggerGroup struct {
	Error log.Logger
	Debug log.Logger
}

type service struct {
	store    WebhookStore
	callback func([]Webhook)
}

func (s *service) Add(owner string, w *Webhook) error {
	fmt.Printf("__Service.go: Add() called with owner %s and webhook %+v\n", owner, w)
	err := s.store.Add(owner, w)
	if err != nil {
		return err
	}
// why passing empty string here
	allWebhooks, err := s.AllWebhooks(owner)
	if err != nil {
		return err
	}

	webhooks := make([]Webhook, len(allWebhooks))
	for i, wh := range allWebhooks {
		webhooks[i] = *wh
	}

	s.callback(webhooks)

	return nil
}

func (s *service) AllWebhooks(owner string) ([]*Webhook, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("__Service.go AllWebhooks Panic :", r)
			debug.PrintStack()
		}
	}()

	webhooks, err := s.store.AllWebhooks(owner)
	if err != nil {
		return nil, err
	}

	webhooksPtr := make([]*Webhook, len(webhooks))

	for i, wh := range webhooks {
		webhooksPtr[i] = wh
	}

	return webhooksPtr, nil
}

func Initialize(cfg *WatchConfig,watches ...Watch) (Service, func(), error) {
	store := NewWebhookStore()

	svc := &service{
		store: store,
		callback: func(webhooks []Webhook) {
			// here watches is empty, so update will never be called
			for _, watch := range watches {
				watch.Update(webhooks)
			}
		},
	}
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

func validateConfig(cfg *WatchConfig) {
	if cfg.WatchUpdateInterval == 0 {
		cfg.WatchUpdateInterval = time.Second * 5
	}
}
