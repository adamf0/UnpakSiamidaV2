package infrastructure

import (
	"context"
	"log"
	"time"
)

func StartOutboxWorker(
	ctx context.Context,
	processor *OutboxProcessor,
) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// log.Println("Outbox processor started")

	for {
		select {
		case <-ctx.Done():
			log.Println("Outbox processor stopped")
			return

		case <-ticker.C:
			if err := processor.Process(ctx); err != nil {
				log.Println("Outbox error:", err)
			}
		}
	}
}
