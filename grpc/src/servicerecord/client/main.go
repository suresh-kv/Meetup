package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "service_proto"

	"google.golang.org/grpc"
)

const (
	address    = "localhost:9001"
	defaultVin = "VIN_TEST"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewServiceHistoryClient(conn)

	// Contact the server and print out its response.
	name := defaultVin
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) //800000
	defer cancel()
	r, err := c.GetLatestServiceRecord(ctx, &pb.Request{Vin: name})
	if err != nil {
		log.Fatalf("could not get details: %v", err)
	}
	t := time.Unix(r.GetTimestamp(), 0)
	strDate := t.Format(time.UnixDate)
	log.Printf("Dealer Name: %s", r.GetDealerName())
	log.Printf("Advisor Name: %s", r.GetAdvisorName())
	log.Printf("Odometer: %d", r.GetOdometer())
	log.Printf("Timestamp: %s", strDate)
	log.Printf("Repair No: %s", r.GetRepairNo())

	log.Printf("\n\n\n\n******Error**********\n\n\n\n")
	r, err = c.AlwaysReturnError(ctx, &pb.Request{Vin: "test"})
	if err != nil {
		log.Printf("Error: %v", err)
	}

	log.Printf("\n\n\n\n******Experiment Streaming**********\n\n\n\n")

	stream, err := c.GetAllServiceRecords(context.Background(), &pb.Request{Vin: name})

	for {
		r, err := stream.Recv()
		if err != nil {
			log.Fatalf("can not receive %v", err)
		}
		log.Printf("Dealer Name: %s", r.GetDealerName())
		log.Printf("Advisor Name: %s", r.GetAdvisorName())
		log.Printf("Odometer: %d", r.GetOdometer())
		log.Printf("Timestamp: %s", strDate)
		log.Printf("Repair No: %s\n\n", r.GetRepairNo())
	}

}
