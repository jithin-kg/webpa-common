package mhook

import (
	"errors"
	"fmt"
	"sync"
)

type WebhookStore interface {
	Add(owner string, w *Webhook) error
	Delete(owner string, url string) error
	AllWebhooks(owner string) ([]*Webhook, error)
}
type webhookStore struct {
	store map[string]map[string]*Webhook // owner -> url -> Webhook
	mu    sync.RWMutex
}

func NewWebhookStore() WebhookStore {
	return &webhookStore{
		store: make(map[string]map[string]*Webhook),
	}
}

func (ws *webhookStore) Add(owner string, w *Webhook) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	// Create the owner's map if it doesn't exist yet
	if ws.store[owner] == nil {
		ws.store[owner] = make(map[string]*Webhook)
	}

	// Check if the webhook already exists
	if _, ok := ws.store[owner][w.Config.URL]; ok {
		return errors.New("webhook already exists")
	}

	// Add the webhook
	ws.store[owner][w.Config.URL] = w
	// Log the added webhook
	fmt.Printf("__WebhookStore.go: Added webhook: %+v\n", *w)
	fmt.Printf("__WebhookStore.go: Add() current store is : %+v\n", ws.store)
	return nil
}

func (ws *webhookStore) Delete(owner string, url string) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	// Check if the owner's map exists
	if ws.store[owner] == nil {
		return errors.New("owner not found")
	}

	// Check if the webhook exists
	if _, ok := ws.store[owner][url]; !ok {
		return errors.New("webhook not found")
	}

	// Delete the webhook
	delete(ws.store[owner], url)

	return nil
}

func (ws *webhookStore) AllWebhooks(owner string) ([]*Webhook, error) {
	fmt.Printf("__WebhookStore.go: AllWebhooks() Owner is: %s\n", owner)
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	// Check if the owner's map exists
	if ws.store[owner] == nil {
		return []*Webhook{}, nil
	}

	// Convert the map of webhooks to a slice
	webhooks := make([]*Webhook, 0, len(ws.store[owner]))
	for _, w := range ws.store[owner] {
		webhooks = append(webhooks, w)
	}
	fmt.Printf("__WebhookStore.go: AllWebhooks() current store is : %+v\n", ws.store)
	return webhooks, nil
}
