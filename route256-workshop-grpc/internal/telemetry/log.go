package telemetry

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func WithUnaryLogs(
	logger *zap.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		fullMethod := info.FullMethod

		resp, err := handler(ctx, req)

		if err != nil {
			logger.Error(
				"error occured while handling the request",
				zap.String("method", fullMethod),
				zap.Error(err),
			)
		}

		logger.Info(
			"served request",
			zap.String("method", fullMethod),
		)

		return resp, err
	}
}
