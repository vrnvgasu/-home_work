package app

import (
	"context"

	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type App struct {
	storage Storage
	l       Logger
}

type Logger interface {
	Warn(msg string)
	Info(msg string)
	Error(msg string)
	File(msg string) // TODO
}

type Storage interface {
	storage.IStorage
}

func New(logger Logger, storage Storage) *App {
	return &App{
		l:       logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, e storage.Event) (storage.Event, error) {
	id, err := a.storage.Add(ctx, e)
	if err != nil {
		a.l.Error(err.Error())
		return storage.Event{}, status.Errorf(codes.Internal, "app CreateEvent Add failed")
	}

	e.ID = id

	return e, nil
}

func (a *App) UpdateEvent(ctx context.Context, e storage.Event) error {
	err := a.storage.Update(ctx, e)
	if err != nil {
		a.l.Error(err.Error())
		return status.Errorf(codes.Internal, "app UpdateEvent Update failed")
	}

	return nil
}

func (a *App) EventList(ctx context.Context, params storage.Params) ([]storage.Event, error) {
	events, err := a.storage.List(ctx, params)
	if err != nil {
		a.l.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "app EventList List failed")
	}

	return events, nil
}
