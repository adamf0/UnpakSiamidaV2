package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"gorm.io/gorm"
)

type OutboxProcessor struct {
	DB         *gorm.DB
	Dispatcher *EventDispatcher
}

func (p *OutboxProcessor) Process(ctx context.Context) error {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[outbox] Process recovered from panic: %v\n", r)
		}
	}()
	if p == nil || p.DB == nil {
		return nil
	}
	var messages []OutboxMessage

	if err := p.DB.
		Where("processed_at IS NULL").
		Limit(10).
		Find(&messages).Error; err != nil {
		return err
	}

	for _, msg := range messages {
		if err := p.processMessage(ctx, &msg); err != nil {
			p.DB.Model(&msg).Update("error", err.Error())
		}
	}

	return nil
}

func (p *OutboxProcessor) processMessage(
	ctx context.Context,
	msg *OutboxMessage,
) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in outbox processing: %v", r)
			log.Printf("[outbox] Recovered panic for %s: %v\n", msg.Type, r)
		}
	}()

	log.Println("[outbox] Execute:", msg.Type)

	eventType := resolveType(msg.Type)
	eventPtr := reflect.New(eventType)

	if err := json.Unmarshal(
		[]byte(msg.Payload),
		eventPtr.Interface(),
	); err != nil {
		return err
	}

	eventValue := eventPtr.Elem().Interface()

	if err := p.Dispatcher.Dispatch(ctx, eventValue); err != nil {
		return err
	}

	now := time.Now().UTC()
	return p.DB.Model(msg).
		Update("processed_at", &now).Error
}
