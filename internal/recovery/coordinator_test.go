package recovery

import (
	"context"
	"errors"
	"testing"
	"time"
)

type restoreStoreStub struct {
	replaceErr error
	replaced   []Record
	events     []RestoreEvent
}

func (s *restoreStoreStub) Replace(_ context.Context, records []Record) error {
	s.replaced = append([]Record(nil), records...)
	return s.replaceErr
}

func (s *restoreStoreStub) AppendEvent(_ context.Context, event RestoreEvent) error {
	s.events = append(s.events, event)
	return nil
}

func TestRestoreReportsPersistenceFailureAfterAuditSucceeds(t *testing.T) {
	now := time.Now().UTC()
	snapshot, err := Build("farm", 1, []Record{{Kind: "animal", ID: "cow-1", TenantID: "farm", Version: 2, UpdatedAt: now, Payload: []byte(`{"status":"active"}`)}}, now)
	if err != nil {
		t.Fatal(err)
	}
	persistenceErr := errors.New("restore storage unavailable")
	store := &restoreStoreStub{replaceErr: persistenceErr}
	err = (Coordinator{Store: store}).Restore(context.Background(), snapshot, nil, nil)
	if !errors.Is(err, persistenceErr) {
		t.Fatalf("restore error = %v, want persistence failure", err)
	}
	if len(store.events) != 1 || store.events[0].Outcome != "failed" || store.events[0].Detail == "" {
		t.Fatalf("events = %+v", store.events)
	}
}
