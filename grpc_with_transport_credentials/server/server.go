package server

import (
	"golang.org/x/net/context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	pb "github.com/takaishi/hello2018/grpc_with_transport_credentials/protocol"
	"github.com/takaishi/hello2018/grpc_with_transport_credentials/sshTC"
	"github.com/urfave/cli"
)

type customerService struct {
	customers []*pb.Person
	m         sync.Mutex
}

func (cs *customerService) ListPerson(p *pb.RequestType, stream pb.CustomerService_ListPersonServer) error {
	cs.m.Lock()
	defer cs.m.Unlock()
	for _, p := range cs.customers {
		if err := stream.Send(p); err != nil {
			return err
		}
	}
	return nil
}

func (cs *customerService) AddPerson(c context.Context, p *pb.Person) (*pb.ResponseType, error) {
	cs.m.Lock()
	defer cs.m.Unlock()
	cs.customers = append(cs.customers, p)
	return new(pb.ResponseType), nil
}

func Start(c *cli.Context) error {
	lis, err := net.Listen("tcp", ":11111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	sshTC := sshTC.NewServerCreds(c.String("public-key"))
	server := grpc.NewServer(grpc.Creds(sshTC))
	pb.RegisterCustomerServiceServer(server, new(customerService))
	return server.Serve(lis)
}
