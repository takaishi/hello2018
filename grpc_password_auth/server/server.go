package server

import (
	"golang.org/x/net/context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	"fmt"
	pb "github.com/takaishi/hello2018/grpc_password_auth/protocol"
	"github.com/takaishi/hello2018/grpc_password_auth/server/auth"
	"github.com/urfave/cli"
)

type customerService struct {
	customers []*pb.Person
	m         sync.Mutex
}

func (cs *customerService) ListPerson(p *pb.RequestType, stream pb.CustomerService_ListPersonServer) error {
	fmt.Println("ListPerson")
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
	fmt.Println("AddPerson")
	cs.m.Lock()
	defer cs.m.Unlock()
	cs.customers = append(cs.customers, p)
	return new(pb.ResponseType), nil
}

func Start(c *cli.Context) error {
	a := auth.NewAuthorizer(c.String("username"), c.String("password"))

	lis, err := net.Listen("tcp", ":11111")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//stm := []grpc.StreamServerInterceptor{}

	server := grpc.NewServer(
		grpc.StreamInterceptor(a.HandleStream),
		grpc.UnaryInterceptor(a.HandleUnary),
	)

	fmt.Printf("server: %#v\n", server)

	pb.RegisterCustomerServiceServer(server, new(customerService))
	return server.Serve(lis)
}
