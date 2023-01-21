package api

import repo2 "gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/repo"
import (
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/model"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/pkg/logger"
	. "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type NotificationEventService struct {
	chanHolder map[uint64]map[string]chan *model.UserNotificationEvent
	sync.RWMutex
}

func NewNotificationEventService() *NotificationEventService {
	return &NotificationEventService{
		chanHolder: map[uint64]map[string]chan *model.UserNotificationEvent{},
	}
}

func (s *NotificationEventService) notify(deviceId uint64, event *model.UserNotificationEvent) {
	for _, session := range s.chanHolder[deviceId] {
		session <- event
	}
}

func (s *NotificationEventService) removeChan(deviceId uint64, key string) {
	s.Lock()
	defer s.Unlock()
	sessions := s.chanHolder[deviceId]
	if sessions != nil {
		delete(sessions, key)
	}
}

func (s *NotificationEventService) registerChan(deviceId uint64, channel chan *model.UserNotificationEvent) string {
	s.Lock()
	defer s.Unlock()
	newUUID, _ := uuid.NewUUID()
	key := newUUID.String()
	sessions := s.chanHolder[deviceId]
	if sessions != nil {
		sessions[key] = channel
		return key
	} else {
		sessions := make(map[string]chan *model.UserNotificationEvent)
		s.chanHolder[deviceId] = sessions
		sessions[key] = channel
	}
	return key
}

type notificationApi struct {
	UnimplementedActNotificationApiServiceServer
	repo         repo2.NotificationRepo
	eventService *NotificationEventService
}

func (a *notificationApi) SendNotificationV1(ctx context.Context, req *SendNotificationV1Request) (*SendNotificationV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.SendNotificationV1")
	defer span.Finish()

	ctx = logger.LogLevelFromContext(ctx)

	if err := req.Validate(); err != nil {
		logger.ErrorKV(
			ctx,
			"SendNotificationV1 -- invalid argument",
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	//now := time.Now()

	event := &model.NotificationEvent{
		DeviceID:  req.Notification.DeviceId,
		Lang:      model.Lang(req.Notification.Lang),
		Message:   req.Notification.Message,
		Status:    model.NotificationCreated,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	err := a.repo.AddNotification(ctx, event)
	if err != nil {
		logger.ErrorKV(
			ctx,
			"SendNotificationV1 -- failed insert record to notification_events table",
			"err", err,
		)
		return nil, status.Error(codes.Internal, err.Error())
	}

	event.BuildUserNotification()
	err = a.repo.UpdatePayload(ctx, event)

	if err != nil {
		logger.ErrorKV(
			ctx,
			"SendNotificationV1 -- failed insert record to notification_events table",
			"err", err,
		)
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer a.eventService.notify(req.Notification.DeviceId, event.UserNotificationEvent)
	response := SendNotificationV1Response{
		NotificationId: event.ID,
	}
	return &response, nil

}

func (a *notificationApi) GetNotification(ctx context.Context, req *GetNotificationV1Request) (*GetNotificationV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.GetNotification")
	defer span.Finish()

	ctx = logger.LogLevelFromContext(ctx)

	if err := req.Validate(); err != nil {
		logger.ErrorKV(
			ctx,
			"GetNotification -- invalid argument",
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	events, err := a.repo.GetNotProcessedByDeviceId(ctx, req.GetDeviceId())
	if err != nil {
		logger.ErrorKV(
			ctx,
			"GetNotification -- invalid argument",
			"err", err,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}
	userEvents := make([]*UserNotification, 0, 10)
	if len(events) > 0 {
		for _, event := range events {
			if event.Status == model.NotificationCreated {
				err := a.repo.UpdateInProgress(ctx, event.ID)
				if err != nil {
					logger.ErrorKV(
						ctx,
						"Update notification to InProgress failed",
						"err", err,
					)
				}
			}

			userEvents = append(userEvents, &UserNotification{
				NotificationId: event.ID,
				Message:        event.UserNotificationEvent.Text,
			})
		}
	}

	return &GetNotificationV1Response{
		Notification: userEvents,
	}, nil

}
func (a *notificationApi) SubscribeNotification(req *SubscribeNotificationRequest, server ActNotificationApiService_SubscribeNotificationServer) error {
	events := make(chan *model.UserNotificationEvent, 128_000)
	key := a.eventService.registerChan(req.DeviceId, events)

	for {
		select {
		case event := <-events:
			{
				err := server.Send(&UserNotification{
					NotificationId: event.NotificationId,
					Message:        event.Text,
				})

				if err != nil {
					logger.ErrorKV(context.Background(), "Closing channel. Remove from holder")
					close(events)
					a.eventService.removeChan(req.DeviceId, key)
					return nil
				}
				err = a.repo.UpdateInProgress(context.Background(), event.NotificationId)

				if err != nil {
					logger.ErrorKV(context.Background(), "Failed to Update notification_event")
				}
			}

		}
	}
	return nil

}
func (a *notificationApi) AckNotification(ctx context.Context, req *AckNotificationV1Request) (*AckNotificationV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.AckNotification")
	defer span.Finish()

	ctx = logger.LogLevelFromContext(ctx)

	if err := req.Validate(); err != nil {
		logger.ErrorKV(
			ctx,
			"GetNotification -- invalid argument",
			"err", err,
		)

		return &AckNotificationV1Response{Success: false}, status.Error(codes.InvalidArgument, err.Error())
	}

	err := a.repo.Ack(ctx, req.NotificationId)
	if err != nil {
		logger.ErrorKV(
			ctx,
			"AckNotification -- update problem",
			"err", err,
		)

		return &AckNotificationV1Response{Success: false}, status.Error(codes.Internal, err.Error())
	}
	return &AckNotificationV1Response{Success: true}, nil
}

// NewNotificationAPI returns api of act-device-api service
func NewNotificationAPI(r repo2.NotificationRepo) ActNotificationApiServiceServer {
	return &notificationApi{
		UnimplementedActNotificationApiServiceServer: UnimplementedActNotificationApiServiceServer{},
		repo:         r,
		eventService: NewNotificationEventService(),
	}
}
