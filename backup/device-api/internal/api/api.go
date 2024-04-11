package api

import (
	"context"
	"database/sql"
	"github.com/opentracing/opentracing-go"
	repo2 "gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/app/repo"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/model"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/pkg/logger"
	. "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	totalDeviceNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Subsystem: "act_device_api",
		Name:      "device_not_found_total",
		Help:      "Total number of devices that were not found",
	})

	cudActionsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "act_device_api",
		Name:      "cud_actions_total",
	}, []string{"action_type"})
)

type deviceAPI struct {
	UnimplementedActDeviceApiServiceServer
	repo      repo2.Repo
	eventRepo repo2.EventRepo
}

// NewDeviceAPI returns api of act-device-api service
func NewDeviceAPI(r repo2.Repo, er repo2.EventRepo) ActDeviceApiServiceServer {
	return &deviceAPI{repo: r, eventRepo: er}
}

func (o *deviceAPI) CreateDeviceV1(
	ctx context.Context,
	req *CreateDeviceV1Request,
) (*CreateDeviceV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.CreateDeviceV1")
	defer span.Finish()

	ctx = logger.LogLevelFromContext(ctx)

	if err := req.Validate(); err != nil {
		logger.ErrorKV(
			ctx,
			"CreateDeviceV1 -- invalid argument",
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	now := time.Now()

	device := &model.Device{
		UserID:    req.GetUserId(),
		Platform:  req.GetPlatform(),
		EnteredAt: &now,
	}

	deviceID, err := o.repo.CreateDevice(ctx, device)
	if err != nil {
		logger.ErrorKV(
			ctx,
			"CreateDeviceV1 -- failed",
			"err", err,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	err = o.eventRepo.Add(ctx, &model.DeviceEvent{
		DeviceID: deviceID,
		Type:     model.Created,
		Status:   model.Deferred,
		Device:   device,
	})
	if err != nil {
		logger.ErrorKV(
			ctx,
			"CreateDeviceV1 -- failed record to event table",
			"err", err,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	cudActionsTotal.WithLabelValues("create").Inc()

	logger.DebugKV(ctx, "CreateDeviceV1 -- success")

	return &CreateDeviceV1Response{
		DeviceId: deviceID,
	}, nil
}

func (o *deviceAPI) DescribeDeviceV1(
	ctx context.Context,
	req *DescribeDeviceV1Request,
) (*DescribeDeviceV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.DescribeDeviceV1")
	defer span.Finish()

	ctx = logger.LogLevelFromContext(ctx)

	if err := req.Validate(); err != nil {
		logger.ErrorKV(
			ctx,
			"DescribeDeviceV1 -- invalid argument",
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	device, err := o.repo.DescribeDevice(ctx, req.GetDeviceId())
	if err != nil && err != sql.ErrNoRows {
		logger.ErrorKV(
			ctx,
			"DescribeDeviceV1 -- failed",
			"err", err,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if device == nil || err == sql.ErrNoRows {
		logger.DebugKV(
			ctx,
			"DescribeDeviceV1 -- device not found",
			"deviceId", req.DeviceId,
		)
		totalDeviceNotFound.Inc()

		return nil, status.Error(codes.NotFound, "device not found")
	}

	logger.DebugKV(ctx, "DescribeDeviceV1 -- success")

	return &DescribeDeviceV1Response{
		Value: &Device{
			Id:        device.ID,
			Platform:  device.Platform,
			UserId:    device.UserID,
			EnteredAt: tspb.New(*device.EnteredAt),
		},
	}, nil
}

func (o *deviceAPI) ListDevicesV1(
	ctx context.Context,
	req *ListDevicesV1Request,
) (*ListDevicesV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.ListDevicesV1")
	defer span.Finish()

	ctx = logger.LogLevelFromContext(ctx)

	if err := req.Validate(); err != nil {
		logger.ErrorKV(
			ctx,
			"ListDevicesV1 -- invalid argument",
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	devices, err := o.repo.ListDevices(ctx, req.GetPage(), req.GetPerPage())
	if err != nil {
		logger.ErrorKV(
			ctx,
			"ListDevicesV1 -- failed",
			"err", err,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	logger.DebugKV(ctx, "ListDevicesV1 -- success")

	var act_device_apiDevices []*Device

	for _, device := range devices {
		act_device_apiDevices = append(act_device_apiDevices,
			&Device{
				Id:        device.ID,
				Platform:  device.Platform,
				UserId:    device.UserID,
				EnteredAt: tspb.New(*device.EnteredAt),
			},
		)
	}

	return &ListDevicesV1Response{
		Items: act_device_apiDevices,
	}, nil
}

func (o *deviceAPI) UpdateDeviceV1(
	ctx context.Context,
	req *UpdateDeviceV1Request,
) (*UpdateDeviceV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.UpdateDeviceV1")
	defer span.Finish()

	ctx = logger.LogLevelFromContext(ctx)

	if err := req.Validate(); err != nil {
		logger.ErrorKV(
			ctx,
			"UpdateDeviceV1 -- invalid argument",
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	deviceID := req.GetDeviceId()

	device := &model.Device{
		ID:       deviceID,
		UserID:   req.GetUserId(),
		Platform: req.GetPlatform(),
	}

	success, err := o.repo.UpdateDevice(ctx, device)
	if err != nil {
		logger.ErrorKV(
			ctx,
			"UpdateDeviceV1 -- failed",
			"err", err,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if success {
		cudActionsTotal.WithLabelValues("update").Inc()

		err = o.eventRepo.Add(ctx, &model.DeviceEvent{
			DeviceID: deviceID,
			Type:     model.Updated,
			Status:   model.Deferred,
			Device:   device,
		})
		if err != nil {
			logger.ErrorKV(
				ctx,
				"UpdateDeviceV1 -- failed record to event table",
				"err", err,
			)

			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	logger.DebugKV(ctx, "UpdateDeviceV1 -- success")

	return &UpdateDeviceV1Response{
		Success: success,
	}, nil
}

func (o *deviceAPI) RemoveDeviceV1(
	ctx context.Context,
	req *RemoveDeviceV1Request,
) (*RemoveDeviceV1Response, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "api.RemoveDeviceV1")
	defer span.Finish()

	ctx = logger.LogLevelFromContext(ctx)

	if err := req.Validate(); err != nil {
		logger.ErrorKV(
			ctx,
			"RemoveDevicesV1 -- invalid argument",
			"err", err,
		)

		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	deviceID := req.GetDeviceId()

	found, err := o.repo.RemoveDevice(ctx, deviceID)
	if err != nil {
		logger.ErrorKV(
			ctx,
			"RemoveDevicesV1 -- failed",
			"err", err,
		)

		return nil, status.Error(codes.Internal, err.Error())
	}

	if !found {
		totalDeviceNotFound.Inc()
	} else {
		cudActionsTotal.WithLabelValues("remove").Inc()

		err = o.eventRepo.Add(ctx, &model.DeviceEvent{
			DeviceID: deviceID,
			Type:     model.Created,
			Status:   model.Deferred,
		})
		if err != nil {
			logger.ErrorKV(
				ctx,
				"RemoveDevicesV1 -- failed record to event table",
				"err", err,
			)

			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	logger.DebugKV(ctx, "RemoveDevicesV1 -- success")

	return &RemoveDeviceV1Response{
		Found: found,
	}, nil
}
