package main

import (
	"net"
	"log"
	"golang.org/x/net/context"
	"github.com/nebulaim/telegramd/baselib/grpc_util/middleware/examples/helloworld"
	"google.golang.org/grpc"
)

// GreeterServer is the server API for Greeter service.
type GreeterServerImpl struct {
}

func (s *GreeterServerImpl) SayHello(ctx context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	r := &helloworld.HelloReply{
		Message: request.Name,
	}
	return r, nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	log.Printf("rpc listening on 0.0.0.0:8100")

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &GreeterServerImpl{})
	// proto.RegisterEchoServiceServer(s.s, s)
	s.Serve(listener)
}
