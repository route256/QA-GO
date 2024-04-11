package repo

import (
	"context"
	"encoding/json"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/model"
)

// EventRepo go:generate mockgen -source=./event.go -destination=./../../mocks/repo_mock.go -package=mocks
type EventRepo interface {
	Lock(ctx context.Context, n uint64) ([]model.DeviceEvent, error)
	Unlock(ctx context.Context, eventIDs []uint64) error

	Add(ctx context.Context, event *model.DeviceEvent) error
	Remove(ctx context.Context, eventIDs []uint64) error
}

// NewEventRepo returns EventRepo interface
func NewEventRepo(db *sqlx.DB, batchSize uint) EventRepo {
	return &repo{db: db, batchSize: batchSize}
}

func (r repo) Lock(ctx context.Context, n uint64) ([]model.DeviceEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.event.Lock")
	defer span.Finish()

	query := sq.Update("devices_events").PlaceholderFormat(sq.Dollar).
		Set("status", model.Processed).
		Set("updated_at", "now()").
		Where(sq.Select("id").PlaceholderFormat(sq.Dollar).
			Prefix("id IN (").
			From("devices_events").
			Where(sq.NotEq{"status": model.Processed}).
			OrderBy("created_at").
			Limit(n).
			Suffix("FOR UPDATE SKIP LOCKED)"),
		).
		// TODO убрал payload. надо добавить корректный парсинг jsonb в объект
		Suffix("RETURNING id, device_id, type, status, created_at, updated_at")

	s, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var events []model.DeviceEvent

	err = r.db.SelectContext(ctx, &events, s, args...)

	return events, err
}

func (r repo) Unlock(ctx context.Context, eventIDs []uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.event.Unlock")
	defer span.Finish()

	query := sq.Update("devices_events").PlaceholderFormat(sq.Dollar).
		Set("status", model.Deferred).
		Set("updated_at", "now()").
		Where(sq.Eq{"device_id": eventIDs})

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, s, args...)

	return err
}

func (r repo) Add(ctx context.Context, event *model.DeviceEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.event.Add")
	defer span.Finish()

	payload, err := json.Marshal(event.Device)
	if err != nil {
		return err
	}

	query := sq.Insert("devices_events").PlaceholderFormat(sq.Dollar).
		Columns("device_id", "type", "status", "payload").
		Values(event.DeviceID, event.Type, event.Status, payload).
		Suffix("RETURNING id")

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, s, args...)

	return err
}

func (r repo) Remove(ctx context.Context, eventIDs []uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.event.Remove")
	defer span.Finish()

	query := sq.Delete("devices_events").PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"device_id": eventIDs})

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, s, args...)

	return err
}
