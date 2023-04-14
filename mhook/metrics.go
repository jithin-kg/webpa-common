package mhook

import (
	"github.com/go-kit/kit/metrics"
	"github.com/jithin-kg/webpa-common/xmetrics"
)

// Names
const (
	PollCounter          = "webhook_polls_total"
	WebhookListSizeGauge = "webhook_list_size_value"
)

// Labels
const (
	OutcomeLabel = "outcome"
)

// Label Values
const (
	SuccessOutcome = "success"
	FailureOutcome = "failure"
)

// Metrics returns the Metrics relevant to this package
func Metrics() []xmetrics.Metric {
	return []xmetrics.Metric{
		{
			Name: WebhookListSizeGauge,
			Type: xmetrics.GaugeType,
			Help: "Size of the current list of webhooks.",
		},
	}
}

type measures struct {
	pollCount       metrics.Counter
	webhookListSize metrics.Gauge
}
