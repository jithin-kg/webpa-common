package mhook

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	// Create a test webhook
	// Create a test webhook
wh := &Webhook{
	Address: "https://client.example.com",
	Config: struct {
			URL             string   `json:"url"`
			ContentType     string   `json:"content_type"`
			Secret          string   `json:"secret,omitempty"`
			AlternativeURLs []string `json:"alt_urls,omitempty"`
	}{
			URL:         "https://example.com/webhook",
			ContentType: "application/json",
			Secret:      "",
	},
	FailureURL: "https://failure.example.com",
	Events:     []string{"event1", "event2"},
	Matcher: struct {
			DeviceID []string `json:"device_id"`
	}{
			DeviceID: []string{"device1", "device2"},
	},
	Duration: time.Second * 5,
	Until:    time.Now().Add(time.Second * 5),
}


	cfg := &WatchConfig{
		Webhooks:           []Webhook{*wh},
		WatchUpdateInterval: time.Second * 5,
	}
// Create a sample function to handle webhook updates
	handleWebhookUpdate := func(webhooks []Webhook) {
		fmt.Println("Webhook update:", webhooks)
	}
	// Create a test watch
	testWatch := WatchFunc(handleWebhookUpdate)
	// Initialize the service
	service, cleanup, err := Initialize(cfg, testWatch)
	require.NoError(t, err)
	defer cleanup()

	// Test adding a webhook
	err = service.Add("owner1", wh)
	assert.NoError(t, err)
// add new webhook


	// Test retrieving webhooks for an owner
	webhooks, err := service.AllWebhooks("owner1")
	require.NoError(t, err)
	assert.Len(t, webhooks, 1)
	assert.Equal(t, wh.Config.URL, webhooks[0].Config.URL)
}
