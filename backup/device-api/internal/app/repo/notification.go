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
type NotificationRepo interface {
	GetNotProcessedByDeviceId(ctx context.Context, n uint64) ([]*model.NotificationEvent, error)

	Ack(ctx context.Context, notificationId uint64) error
	UpdateInProgress(ctx context.Context, notificationId uint64) error
	AddNotification(ctx context.Context, event *model.NotificationEvent) error
	UpdatePayload(ctx context.Context, event *model.NotificationEvent) error
}

// NewEventRepo returns EventRepo interface
func NewNotificationRepo(db *sqlx.DB, batchSize uint) NotificationRepo {
	return &repo{db: db, batchSize: batchSize}
}
func (r *repo) UpdatePayload(ctx context.Context, event *model.NotificationEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.notification.event.UpdatePayload")
	defer span.Finish()
	notificationEvent := event.UserNotificationEvent
	marshal, err := json.Marshal(notificationEvent)
	if err != nil {
		return err
	}

	query := sq.Update("notification_events").PlaceholderFormat(sq.Dollar).
		Set("payload", marshal).
		Where(sq.Eq{"id": event.ID})

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, s, args...)
	if err != nil {
		return err
	}
	return nil

}

func (r *repo) UpdateInProgress(ctx context.Context, notificationId uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.notification.event.UpdateNotificationToInProgress")
	defer span.Finish()

	query := sq.Update("notification_events").PlaceholderFormat(sq.Dollar).
		Set("status", model.NotificationProcessed).
		Where(sq.Eq{"id": notificationId})

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, s, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) AddNotification(ctx context.Context, event *model.NotificationEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.notification.event.Add")
	defer span.Finish()

	payload, err := json.Marshal(event.UserNotificationEvent)
	if err != nil {
		return err
	}

	query := sq.Insert("notification_events").PlaceholderFormat(sq.Dollar).
		Columns("device_id", "message", "status", "payload", "lang").
		Values(event.DeviceID, event.Message, event.Status, payload, event.Lang).
		Suffix("RETURNING id, created_at, updated_at")

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}
	var retEvent model.NotificationEvent
	err = r.db.GetContext(ctx, &retEvent, s, args...)
	event.ID = retEvent.ID
	event.CreatedAt = retEvent.CreatedAt
	event.UpdatedAt = retEvent.UpdatedAt

	return err
}

func (r *repo) GetNotProcessedByDeviceId(ctx context.Context, device_id uint64) ([]*model.NotificationEvent, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.notification.event.GetByDeviceId")
	defer span.Finish()
	query := sq.Select("*").PlaceholderFormat(sq.Dollar).
		From("notification_events").
		Where(sq.And{sq.Eq{"device_id": device_id}, sq.Or{sq.Eq{"status": model.NotificationCreated}, sq.Eq{"status": model.NotificationProcessed}}}).
		OrderBy("created_at DESC")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var events []*model.NotificationEvent
	err = r.db.SelectContext(ctx, &events, sql, args...)
	return events, err
}

func (r *repo) Ack(ctx context.Context, notificationId uint64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.notification.event.Ack")
	defer span.Finish()

	query := sq.Update("notification_events").PlaceholderFormat(sq.Dollar).
		Set("status", model.NotificationDelivered).
		Where(sq.Eq{"id": notificationId})

	s, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, s, args...)
	if err != nil {
		return err
	}
	return nil
}
