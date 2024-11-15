package storage

import (
	"context"
	"time"
)

type Params struct {
	Limit      int
	StartAtGEq time.Time `db:"start_at_g_eq"`
}

type IStorage interface {
	Add(context.Context, Event) (uint64, error)
	Update(context.Context, Event) error
	Delete(ctx context.Context, eventID uint64) error
	List(context.Context, Params) ([]Event, error)
	Connect(context.Context) error
	Close(context.Context) error
	Migrate() error
}
