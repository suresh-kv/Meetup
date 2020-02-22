package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	pb "service_proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	port = ":9001"
)

// server is used to implement server_proto.GetLatestServiceRecord
type server struct {
	repairNo int64
}

// GetLatestServiceRecord implements server_proto.GetLatestServiceRecord
func (s *server) GetLatestServiceRecord(ctx context.Context, in *pb.Request) (*pb.ServiceRecord, error) {
	log.Printf("Received: %v", in.GetVin())
	time1 := time.Now().Unix()

	return &pb.ServiceRecord{DealerName: "Dealer",
		Timestamp:   time1,
		Odometer:    30000,
		RepairNo:    "REPAIRNO_1000",
		AdvisorName: "Sam",
	}, nil
}

// GetLatestServiceRecord implements server_proto.GetLatestServiceRecord
func (s *server) GetAllServiceRecords(req *pb.Request, srv pb.ServiceHistory_GetAllServiceRecordsServer) error {
	log.Printf("Received: %v", req.GetVin())
	for i := 0; i < 3; i++ {
		time1 := time.Now().Unix()
		s.repairNo += 1000
		resp := &pb.ServiceRecord{DealerName: "Dealer1",
			Timestamp:   time1,
			Odometer:    30000 + s.repairNo,
			RepairNo:    "REPAIRNO_" + strconv.FormatInt(s.repairNo, 10),
			AdvisorName: "Sam",
		}
		if err := srv.Send(resp); err != nil {
			log.Printf("send error %v", err)
		}
	}
	return nil
}

// GetLatestServiceRecord implements server_proto.GetLatestServiceRecord
func (s *server) AlwaysReturnError(ctx context.Context, in *pb.Request) (*pb.ServiceRecord, error) {
	return nil, status.Errorf(
		codes.PermissionDenied,
		fmt.Sprintf("Your message %v", in.GetVin()),
	)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterServiceHistoryServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)

	}
}
