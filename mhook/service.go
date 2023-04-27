package mhook

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"io/ioutil"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/metrics/generic"

	"gopkg.in/yaml.v2"
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
	logger   *loggerGroup
}

func (s *service) Add(owner string, w *Webhook) error {
	s.logger.Debug.Log("msg", "Add() called", "owner", owner, "webhook", fmt.Sprintf("%+v", w))
	err := s.store.Add(owner, w)
	if err != nil {
		return err
	}
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
	s.logger.Debug.Log("msg", "AllWebhooks called", "owner", owner)
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error.Log("msg", "AllWebhooks Panic", "panic", r)
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

func Initialize(cfg *WatchConfig, watches ...Watch) (Service, func(), error) {
	rootLogger := log.NewLogfmtLogger(os.Stdout)
	rootLogger = log.With(rootLogger, "ts", log.DefaultTimestampUTC)
	rootLogger = log.With(rootLogger, "caller", log.DefaultCaller)
	loggers := newLoggerGroup(rootLogger)
	loggers.Debug.Log("msg", "initialize called")
	watches = append(watches, webhookListSizeWatch(generic.NewGauge(WebhookListSizeGauge)))
	store := NewWebhookStore(loggers)

	svc := &service{
		store:  store,
		logger: loggers,
		callback: func(webhooks []Webhook) {
			// here watches is empty, so update will never be called
			for _, watch := range watches {
				watch.Update(webhooks)
			}
		},
	}
	return svc, func() { /*...*/ }, nil
}
func (s *service) AddWebhookFromYaml(yamlFile string) error {
	data, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return err
	}

	var webhook Webhook
	err = yaml.Unmarshal(data, &webhook)
	if err != nil {
		return err
	}

	s.logger.Debug.Log("msg", "Webhook read from YAML file", "webhook", fmt.Sprintf("%+v", webhook))

	err = s.Add("yaml_owner", &webhook)
	if err != nil {
		return err
	}
	return nil
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
