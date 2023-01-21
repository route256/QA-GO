package main

import (
	"context"
	"log"
	"time"

	"gitlab.ozon.dev/svpetrov/route256-workshop-grpc/pkg/api/dns"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:8888",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := dns.NewDNSClient(conn)

	timeoutCtx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*5,
	)
	defer cancel()

	resp, err := client.GetAddress(timeoutCtx, &dns.DNSService{
		Name: "dummy",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}
