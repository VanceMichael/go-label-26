package recovery

import (
	"context"
	"fmt"

	"go-base/internal/domain"
)

type RestoreStore interface {
	Replace(context.Context, []Record) error
	AppendEvent(context.Context, RestoreEvent) error
}

type Coordinator struct {
	Store RestoreStore
}

func (c Coordinator) Restore(ctx context.Context, snapshot Snapshot, current []Record, allowedKinds map[string]bool) error {
	if c.Store == nil {
		return fmt.Errorf("%w: restore store", domain.ErrInvalid)
	}
	plan, err := Plan(snapshot, current, allowedKinds)
	if err != nil {
		return err
	}
	restored, err := Apply(current, plan)
	if err != nil {
		return err
	}
	event := RestoreEvent{TenantID: snapshot.TenantID, Inserted: len(plan.Insert), Updated: len(plan.Update), Outcome: "completed"}
	if err := c.Store.Replace(ctx, restored); err != nil {
		event.Outcome = "failed"
		event.Detail = err.Error()
		return c.Store.AppendEvent(ctx, event)
	}
	return c.Store.AppendEvent(ctx, event)
}
