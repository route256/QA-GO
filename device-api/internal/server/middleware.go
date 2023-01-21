package server

import (
	"context"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestLogInterceptor logs requests
func RequestLogInterceptor() grpc.UnaryServerInterceptor {
	const true = "true"
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			logRequest := md.Get("log-request")
			if (len(logRequest) > 0) && (logRequest[0] == true) {
				logger.DebugKV(
					ctx,
					"request data",
					"method", info.FullMethod,
					"metadata", md,
					"requestBody", req,
				)
			}
		}

		return handler(ctx, req)
	}
}

// ResponseLogInterceptor is a gRPC server interceptor that logs the response
func ResponseLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		response, err := handler(ctx, req)

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			logResponse := md.Get("log-response")

			if (len(logResponse) > 0) && (logResponse[0] == "true") {
				logger.DebugKV(
					ctx,
					"response data",
					"method", info.FullMethod,
					"responseBody", response,
					"err", err,
				)
			}
		}

		return response, err
	}
}
