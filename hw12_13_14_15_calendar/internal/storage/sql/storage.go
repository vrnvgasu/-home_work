package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/stdlib" // postgresql provider
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/vrnvgasu/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sqlx.DB
}

func New() storage.IStorage {
	return &Storage{}
}

func (s *Storage) Add(ctx context.Context, e storage.Event) (uint64, error) {
	q := `insert into event (title, start_at, end_at, description, owner_id, send_before)
	values ($1, $2, $3, $4, $5, $6) returning id;`
	err := s.db.QueryRowxContext(ctx, q, e.Title, e.StartAt, e.EndAt, e.Description, e.OwnerID, e.SendBefore).
		Scan(&e.ID)
	if err != nil {
		return 0, fmt.Errorf("insert event error: %w", err)
	}

	return e.ID, nil
}

func (s *Storage) Update(ctx context.Context, e storage.Event) error {
	q := `update event
	set title = :title, start_at = :start_at,
    end_at = :end_at, description = :description, 
    owner_id = :owner_id, send_before = :send_before
	where id = :id;`
	_, err := s.db.NamedExecContext(ctx, q, e)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("id not found, %d: %w", e.ID, err)
	} else if err != nil {
		return fmt.Errorf("update event error: %w", err)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, eventID uint64) error {
	q := `delete from event where id = $1`
	_, err := s.db.ExecContext(ctx, q, eventID)
	if err != nil {
		return fmt.Errorf("delete event error: %w", err)
	}

	return nil
}

func (s *Storage) List(ctx context.Context, params storage.Params) ([]storage.Event, error) {
	q := `select * from event`

	if !params.StartAtGEq.IsZero() {
		q += ` where start_at >= :start_at_g_eq `
	}

	if !params.StartAtLEq.IsZero() {
		if !params.StartAtGEq.IsZero() {
			q += ` and `
		}
		q += ` start_at <= :start_at_l_eq `
	}

	q += " order by start_at"

	if params.Limit > 0 {
		q += fmt.Sprintf(" limit %d", params.Limit)
	}

	rows, err := s.db.NamedQueryContext(ctx, q, params)
	if err != nil {
		return nil, fmt.Errorf("list event error: %w", err)
	}
	defer rows.Close()

	result := []storage.Event{}
	for rows.Next() {
		e := storage.Event{}
		if err = rows.StructScan(&e); err != nil {
			return nil, fmt.Errorf("list event scan error: %w", err)
		}

		result = append(result, e)
	}

	return result, nil
}

func (s *Storage) Connect(_ context.Context) error {
	db, err := sqlx.Open("pgx", config.Cfg.PSQL.DSN)
	if err != nil {
		return fmt.Errorf("open sql db fail: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("ping sql db fail: %w", err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) Migrate() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set dialect failed: %w", err)
	}

	if err := goose.Up(s.db.DB, config.Cfg.PSQL.Migration); err != nil {
		return fmt.Errorf("up migration failed: %w", err)
	}

	return nil
}
