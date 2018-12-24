package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/OlgaDnepr/adapter/pb"
	"google.golang.org/grpc"
)

// TODO: move it to config
const adapterURL = "adapter:50001"

func messageRequest(isMarco bool) pb.MarcoPolo {
	if isMarco {
		return pb.MarcoPolo_Marco
	}
	return pb.MarcoPolo_Polo
}

func main() {
	marcoFlag := flag.Bool("marco", true, "A parameter for a client to send to a server.\nIt sends 'Marco' if true or 'Polo' otherwise")
	flag.Parse()

	conn, err := grpc.Dial(adapterURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("client: failed to connect to %s. Err: %s", adapterURL, err)
	}
	defer func() {
		if err = conn.Close(); err != nil {
			log.Fatalf("client: failed to close connection. Err: %s", err)
		}
	}()

	var (
		adapterClient = pb.NewAdapterClient(conn)
		sendRequest   = &pb.Request{
			Message: messageRequest(*marcoFlag),
		}
	)

	gotRequest, err := adapterClient.Get(context.Background(), sendRequest)
	if err != nil {
		log.Fatalf("client: error in adapter communication: %s", err)
	}
	if gotRequest == nil {
		log.Fatalf("client: no reply from adapter")
	}

	fmt.Println(gotRequest.Message)
}
