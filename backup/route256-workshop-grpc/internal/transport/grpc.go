package transport

import (
	context "context"

	"gitlab.ozon.dev/svpetrov/route256-workshop-grpc/pkg/api/dns"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type GRPCDNS struct {
	dns.UnimplementedDNSServer

	implementation DNSImplementation
}

type DNSImplementation interface {
	Register(
		ctx context.Context,
		serviceName string,
		addresses ...string,
	)
	Unregister(
		ctx context.Context,
		serviceName string,
		addresses ...string,
	)
	GetAddressesForService(
		ctx context.Context,
		serviceName string,
	) ([]string, bool)
}

// GetAddress implements dns.DNSServer
func (g GRPCDNS) GetAddress(
	ctx context.Context,
	service *dns.DNSService,
) (*dns.DNSServiceToAddressBinding, error) {
	addresses, _ := g.implementation.GetAddressesForService(
		ctx, service.GetName(),
	)

	resp := &dns.DNSServiceToAddressBinding{
		Service: service,
		Address: []*dns.DNSAddress{},
	}

	for _, v := range addresses {
		resp.Address = append(resp.Address, &dns.DNSAddress{
			Address: v,
		})
	}

	return resp, nil
}

// Register implements dns.DNSServer
func (g GRPCDNS) Register(
	ctx context.Context,
	req *dns.DNSServiceToAddressBinding,
) (*emptypb.Empty, error) {

	g.implementation.Register(
		ctx,
		req.GetService().GetName(),
		pbSliceToStringSlice(req.GetAddress()...)...,
	)

	return &emptypb.Empty{}, nil
}

// Unregister implements dns.DNSServer
func (g GRPCDNS) Unregister(
	ctx context.Context,
	req *dns.DNSServiceToAddressBinding,
) (*emptypb.Empty, error) {

	g.implementation.Unregister(
		ctx, req.GetService().GetName(), pbSliceToStringSlice(req.GetAddress()...)...,
	)

	return &emptypb.Empty{}, nil
}

func NewGRPC(
	implemenetation DNSImplementation,
) dns.DNSServer {
	return GRPCDNS{
		implementation: implemenetation,
	}
}
