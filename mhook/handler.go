package mhook

import (
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
)

func NewAddWebhookHandler(s Service) http.Handler {
	fmt.Println("__hanlder.go: NewAddWebhookHandler() called")
	return kithttp.NewServer(
		newAddWebhookEndpoint(s),
		decodeAddWebhookRequest,
		encodeAddWebhookResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)
}

func NewGetAllWebhooksHandler(s Service) http.Handler {
	return kithttp.NewServer(
		newGetAllWebhooksEndpoint(s),
		decodeGetAllWebhooksRequest,
		encodeGetAllWebhooksResponse,
		kithttp.ServerErrorEncoder(errorEncoder),
	)
}
