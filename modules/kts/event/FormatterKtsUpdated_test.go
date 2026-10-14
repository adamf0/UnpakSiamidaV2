package event

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRenderKtsUpdatedTemplate_NilDates(t *testing.T) {
	event := KtsUpdatedEvent{
		EventID:             uuid.New(),
		KtsUUID:             uuid.New(),
		OccurredOn:          time.Now(),
		TanggalClosing:      nil,
		TanggalPenyelesaian: nil,
		TanggalClosingFinal: nil,
	}

	result := RenderKtsUpdatedTemplate(event)
	if !strings.Contains(result, "Perubahan Data KTS") {
		t.Errorf("Expected result to contain 'Perubahan Data KTS', got: %s", result)
	}
}

func TestRenderKtsUpdatedTemplate_WithDates(t *testing.T) {
	tglClosing := "2026-10-07"
	tglSelesai := "2026-10-08"
	tglFinal := "2026-10-09"

	event := KtsUpdatedEvent{
		EventID:             uuid.New(),
		KtsUUID:             uuid.New(),
		OccurredOn:          time.Now(),
		TanggalClosing:      &tglClosing,
		TanggalPenyelesaian: &tglSelesai,
		TanggalClosingFinal: &tglFinal,
	}

	result := RenderKtsUpdatedTemplate(event)
	if !strings.Contains(result, "07 Oktober 2026") {
		t.Errorf("Expected result to contain formatted date '07 Oktober 2026', got: %s", result)
	}
}
