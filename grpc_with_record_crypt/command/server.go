package command

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/takaishi/hello2018/grpc_with_record_crypt/protocol"
	"github.com/takaishi/hello2018/grpc_with_record_crypt/tc"
	"github.com/urfave/cli"
	"sync"
)

type helloService struct {
	m sync.Mutex
}

func (hs *helloService) Send(c context.Context, req *pb.Request) (*pb.Response, error) {
	hs.m.Lock()
	defer hs.m.Unlock()

	log.Printf("[DEBUG] Request: %+v\n", req)

	resp := &pb.Response{Msg: "Thanks!!!"}
	return resp, nil
}

func StartServer(c *cli.Context, secure bool) error {
	lis, err := net.Listen("tcp", "127.0.0.1:11111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tc := tc.NewServerCreds()
	server := grpc.NewServer(grpc.Creds(tc))

	pb.RegisterHelloServiceServer(server, new(helloService))
	return server.Serve(lis)
}
