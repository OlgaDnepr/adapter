package main

import (
	"context"
	"log"
	"net"

	"github.com/OlgaDnepr/adapter/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: move it to config
const port = ":50002"

// server is used to implement pb.ServerServer
type server struct{}

// Get implements pb.ServerServer
func (*server) Get(ctx context.Context, in *pb.Reply) (*pb.Reply, error) {
	if in == nil {
		return nil, errors.New("server: no request")
	}

	switch in.Message {
	case pb.MonkeyFollow_Follow:
		return &pb.Reply{Message: pb.MonkeyFollow_Monkey}, nil
	case pb.MonkeyFollow_Monkey:
		return &pb.Reply{Message: pb.MonkeyFollow_Follow}, nil
	}

	return nil, errors.Errorf("server: incorrect message %q", in.Message)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("server: failed to listen: %s", err)
	}

	s := grpc.NewServer()
	pb.RegisterServerServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("server: failed to serve: %s", err)
	}
}
