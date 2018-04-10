package client

import (
	"fmt"
	"io"
	pb "github.com/takaishi/hello2018/grpc_password_auth/protocol"
	sshTC2 "github.com/takaishi/hello2018/grpc_with_transport_credentials/sshTC"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/urfave/cli"
)

func Add(c *cli.Context, name string, age int) error {
	sshTC := sshTC2.NewClientCreds()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(sshTC),
		//grpc.WithBlock(),
	}
	conn, err := grpc.Dial("127.0.0.1:11111", opts...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	client := pb.NewCustomerServiceClient(conn)
	person := &pb.Person{
		Name: name,
		Age:  int32(age),
	}
	_, err = client.AddPerson(context.Background(), person)
	return err
}

func List(c *cli.Context) error {
	sshTC := sshTC2.NewClientCreds()
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(sshTC),
		//grpc.WithBlock(),
	}

	conn, err := grpc.Dial("127.0.0.1:11111", opts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewCustomerServiceClient(conn)

	stream, err := client.ListPerson(context.Background(), new(pb.RequestType))
	if err != nil {
		return err
	}
	for {
		person, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(person)
	}
	return nil
}
