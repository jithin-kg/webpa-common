package mhook

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

// Watch is the interface for listening for webhook subcription updates.
// Updates represent the latest known list of subscriptions.
type Watch interface {
	Update([]Webhook)
}

// WatchFunc allows bare functions to pass as Watches.
type WatchFunc func([]Webhook)

func (f WatchFunc) Update(webhooks []Webhook) {
	fmt.Println("__watch.go Update() called")
	f(webhooks)
}

// Config provides the different options for the initializing the wehbook service.
type WatchConfig struct {
	// Webhooks contains the list of webhooks to be used by the webhook service
	Webhooks []Webhook

	// WatchUpdateInterval is the duration between each update to all watchers.
	WatchUpdateInterval time.Duration
}
/**
 we use the webhookListSizeWatch function to create a Watch that is specifically designed to update a metrics.Gauge. 
 This function takes a metrics.Gauge as a parameter and returns a Watch that updates that metrics.
 Gauge with the size of the webhook list on each update. By using this approach, we can separate the concerns of updating the metrics.
 Gauge and updating the webhook list, and allow for greater flexibility in how we handle updates to the webhook list.
**/
func webhookListSizeWatch(s metrics.Gauge) Watch {
	return WatchFunc(func(webhooks []Webhook) {
		s.Set(float64(len(webhooks)))
	})
}
