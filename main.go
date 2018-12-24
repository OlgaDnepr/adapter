package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/OlgaDnepr/adapter/pb"
	"github.com/pkg/errors"
)

// TODO: move it to config
const (
	port      = ":50001"
	serverURL = "server:50002"
)

// server is used to implement pb.AdapterServer
type adapter struct {
	serverClient pb.ServerClient
}

func newAdapterServer(c pb.ServerClient) pb.AdapterServer {
	return &adapter{serverClient: c}
}

// Get implements pb.AdapterServer
func (a *adapter) Get(ctx context.Context, in *pb.Request) (*pb.Request, error) {
	if in == nil {
		return nil, errors.New("adapter: no request")
	}
	replyMessage, err := translateRequest(in.Message)
	if err != nil {
		return nil, errors.Wrap(err, "adapter: ")
	}

	var reply = &pb.Reply{
		Message: replyMessage,
	}
	out, err := a.serverClient.Get(ctx, reply)
	if err != nil {
		return nil, errors.Wrap(err, "adapter: error in server communication")
	}
	if out == nil {
		return nil, errors.Errorf("adapter: no reply from server")
	}

	requestMessage, err := translateReply(out.Message)
	if err != nil {
		return nil, errors.Wrap(err, "adapter: ")
	}

	return &pb.Request{Message: requestMessage}, nil
}

func translateRequest(in pb.MarcoPolo) (pb.MonkeyFollow, error) {
	switch in {
	case pb.MarcoPolo_Marco:
		return pb.MonkeyFollow_Monkey, nil
	case pb.MarcoPolo_Polo:
		return pb.MonkeyFollow_Follow, nil
	}
	return pb.MonkeyFollow(-1), errors.Errorf("incorrect request message %q", in)
}

func translateReply(in pb.MonkeyFollow) (pb.MarcoPolo, error) {
	switch in {
	case pb.MonkeyFollow_Monkey:
		return pb.MarcoPolo_Marco, nil
	case pb.MonkeyFollow_Follow:
		return pb.MarcoPolo_Polo, nil
	}
	return pb.MarcoPolo(-1), errors.Errorf("incorrect reply message %q", in)
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("adapter: failed to listen: %s", err)
	}

	conn, err := grpc.Dial(serverURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("adapter: failed to connect to %s. Err: %s", serverURL, err)
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Fatalf("adapter: failed to close connection. Err: %s", err)
		}
	}()

	var (
		server       = grpc.NewServer()
		serverClient = pb.NewServerClient(conn)
	)

	pb.RegisterAdapterServer(server, newAdapterServer(serverClient))
	if err := server.Serve(lis); err != nil {
		log.Fatalf("adapter: failed to serve: %s", err)
	}
}
