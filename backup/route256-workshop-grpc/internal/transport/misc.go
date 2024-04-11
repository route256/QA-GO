package transport

import "gitlab.ozon.dev/svpetrov/route256-workshop-grpc/pkg/api/dns"

func pbSliceToStringSlice(in ...*dns.DNSAddress) []string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = v.GetAddress()
	}

	return out
}
