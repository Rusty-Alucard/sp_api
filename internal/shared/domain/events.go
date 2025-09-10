package domain

import (
	"context"
	"fmt"
	"sync"
)

// EventHandler represents a domain event handler
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
	EventType() string
}

// EventBus represents a domain event bus
type EventBus interface {
	Publish(ctx context.Context, events ...DomainEvent) error
	Subscribe(handler EventHandler)
	Unsubscribe(eventType string, handler EventHandler)
}

// LocalEventBus implements an in-memory event bus
type LocalEventBus struct {
	handlers map[string][]EventHandler
	mutex    sync.RWMutex
}

// NewLocalEventBus creates a new local event bus
func NewLocalEventBus() *LocalEventBus {
	return &LocalEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Publish publishes domain events
func (bus *LocalEventBus) Publish(ctx context.Context, events ...DomainEvent) error {
	bus.mutex.RLock()
	defer bus.mutex.RUnlock()

	for _, event := range events {
		eventType := event.EventType()
		handlers, exists := bus.handlers[eventType]
		if !exists {
			continue
		}

		for _, handler := range handlers {
			if err := handler.Handle(ctx, event); err != nil {
				return fmt.Errorf("failed to handle event %s: %w", eventType, err)
			}
		}
	}

	return nil
}

// Subscribe subscribes to domain events
func (bus *LocalEventBus) Subscribe(handler EventHandler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	eventType := handler.EventType()
	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
}

// Unsubscribe unsubscribes from domain events
func (bus *LocalEventBus) Unsubscribe(eventType string, handler EventHandler) {
	bus.mutex.Lock()
	defer bus.mutex.Unlock()

	handlers, exists := bus.handlers[eventType]
	if !exists {
		return
	}

	for i, h := range handlers {
		if h == handler {
			bus.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// EventDispatcher handles domain event dispatching
type EventDispatcher struct {
	eventBus EventBus
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher(eventBus EventBus) *EventDispatcher {
	return &EventDispatcher{
		eventBus: eventBus,
	}
}

// Dispatch dispatches domain events
func (d *EventDispatcher) Dispatch(ctx context.Context, events ...DomainEvent) error {
	return d.eventBus.Publish(ctx, events...)
}

