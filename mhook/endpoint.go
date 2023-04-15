package mhook

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
)

func newAddWebhookEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*addWebhookRequest)
		fmt.Println("__AddWebhookEndpoint.go: newAddWebhookEndpoint called with owner:", r.owner, "and webhook:", r.webhook)
		return nil, s.Add(r.owner, r.webhook)
	}
}

func newGetAllWebhooksEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(*getAllWebhooksRequest)
		return s.AllWebhooks(r.owner)
	}
}
