package memorystorage

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	idSequence atomic.Uint64
	events     map[uint64]storage.Event
	mu         sync.RWMutex
}

func New() storage.IStorage {
	return &Storage{
		events: make(map[uint64]storage.Event),
	}
}

func (s *Storage) Add(_ context.Context, event storage.Event) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.idSequence.Add(1)
	event.ID = s.idSequence.Load()
	s.events[event.ID] = event

	return s.idSequence.Load(), nil
}

func (s *Storage) Update(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return fmt.Errorf("event not found, %d: %w", event.ID, storage.ErrNotFound)
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) Delete(_ context.Context, eventID uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, eventID)

	return nil
}

func (s *Storage) List(_ context.Context, params storage.Params) ([]storage.Event, error) {
	result := make([]storage.Event, 0, len(s.events))
	for _, event := range s.events {
		if event.StartAt.Before(params.StartAtGEq) {
			continue
		}

		result = append(result, event)
	}

	slices.SortFunc(result, func(a, b storage.Event) int {
		return a.StartAt.Compare(b.StartAt)
	})

	if params.Limit > 0 {
		limit := params.Limit
		if limit > len(result) {
			limit = len(result)
		}

		result = result[:limit]
	}

	return result, nil
}

func (s *Storage) Connect(_ context.Context) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}

func (s *Storage) Migrate() error {
	return nil
}
