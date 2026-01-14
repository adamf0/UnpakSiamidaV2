package infrastructure

import (
	"context"
	"log"

	"github.com/mehdihadeli/go-mediatr"
)

type EventDispatcher struct {
	handlers map[string]func(context.Context, any) error
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string]func(context.Context, any) error),
	}
}

func RegisterEvent[T any](d *EventDispatcher) {
	var zero T
	key := CanonicalTypeName(zero)

	d.handlers[key] = func(ctx context.Context, e any) error {
		return mediatr.Publish[T](ctx, e.(T))
	}
}

func (d *EventDispatcher) Dispatch(ctx context.Context, event any) error {
	key := CanonicalTypeName(event)

	if h, ok := d.handlers[key]; ok {
		return h(ctx, event)
	}

	log.Printf("⚠️ no event handler registered for %s\n", key)
	return nil
}
