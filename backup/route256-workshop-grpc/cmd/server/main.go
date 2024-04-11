package main

import (
	"log"
	"net"

	"gitlab.ozon.dev/svpetrov/route256-workshop-grpc/internal/service"
	"gitlab.ozon.dev/svpetrov/route256-workshop-grpc/internal/telemetry"
	"gitlab.ozon.dev/svpetrov/route256-workshop-grpc/internal/transport"
	"gitlab.ozon.dev/svpetrov/route256-workshop-grpc/pkg/api/dns"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(telemetry.WithUnaryLogs(logger)),
		grpc.Creds(insecure.NewCredentials()))

	dnsProvider := service.NewInMemDNS()

	grpcTransport := transport.NewGRPC(dnsProvider)
	dns.RegisterDNSServer(grpcServer, grpcTransport)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve gRPC server", zap.Error(err))
	}
}

// func withLog() grpc.ServerOption {
// 	return grpc.UnaryInterceptor(
// 		func(
// 			ctx context.Context,
// 			req interface{},
// 			info *grpc.UnaryServerInfo,
// 			handler grpc.UnaryHandler) (resp interface{}, err error) {

// 			resp, err = handler(ctx, req)

// 			return resp, err
// 		})
// }
