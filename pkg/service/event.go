package service

import (
	"github.com/ncarlier/readflow/pkg/event"
	"github.com/ncarlier/readflow/pkg/event/dispatcher"
	"github.com/ncarlier/readflow/pkg/metric"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/rs/zerolog/log"
)

const (
	// EventCreateUser is the create event on an user
	EventCreateUser = "user:create"
	// EventUpdateUser is the update event on an user
	EventUpdateUser = "user:update"
	// EventDeleteUser is the delete event on an user
	EventDeleteUser = "user:delete"
	// EventCreateArticle is the create event on an article
	EventCreateArticle = "article:create"
	// EventUpdateArticle is the update event on an article
	EventUpdateArticle = "article:update"
)

const (
	// NoNotificationEventOption will disable global notification policy
	NoNotificationEventOption event.EventOption = 1 << iota // 1
)

func (reg *Registry) registerEventHandlers() {
	if reg.dispatcher != nil {
		reg.events.Subscribe(EventCreateUser, newExternalEventHandler(reg.dispatcher))
		reg.events.Subscribe(EventUpdateUser, newExternalEventHandler(reg.dispatcher))
		reg.events.Subscribe(EventDeleteUser, newExternalEventHandler(reg.dispatcher))
	}
	reg.events.Subscribe(EventCreateArticle, newCreateArticleMetricEventHandler())
	reg.events.Subscribe(EventCreateArticle, newNotificationEventHandler(reg))
}

func newExternalEventHandler(dis dispatcher.Dispatcher) event.EventHandler {
	return func(evt event.Event) {
		externalEvent := dispatcher.NewExternalEvent(evt.Name, evt.Payload)
		if err := dis.Dispatch(externalEvent); err != nil {
			log.Error().Err(err).Str("event", evt.Name).Msg("unable to send external event")
		}
	}
}

func newCreateArticleMetricEventHandler() event.EventHandler {
	return func(evt event.Event) {
		if article, ok := evt.Payload.(model.Article); ok {
			metric.IncNewArticlesCounter(article)
		}
	}
}
