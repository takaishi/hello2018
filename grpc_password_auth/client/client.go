package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/mattn/sc"
	pb "github.com/takaishi/hello2018/grpc_password_auth/protocol"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type loginCreds struct {
	Username, Password string
}

func (c *loginCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	fmt.Println("GetRequestMetadata")
	return map[string]string{
		"username": c.Username,
		"password": c.Password,
	}, nil
}

func (c *loginCreds) RequireTransportSecurity() bool {
	return false
}

func add(name string, age int) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&loginCreds{Username: "admin", Password: "admin123"}),
	}
	conn, err := grpc.Dial("127.0.0.1:11111", opts...)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	client := pb.NewCustomerServiceClient(conn)
	fmt.Printf("client: %#v\n", client)

	person := &pb.Person{
		Name: name,
		Age:  int32(age),
	}
	_, err = client.AddPerson(context.Background(), person)
	return err
}

func list() error {
	conn, err := grpc.Dial("127.0.0.1:11111",
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&loginCreds{
			Username: "admin",
			Password: "admin123a",
		},
		))
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

func main() {
	(&sc.Cmds{
		{
			Name: "list",
			Desc: "list: listing person",
			Run: func(c *sc.C, args []string) error {
				return list()
			},
		},
		{
			Name: "add",
			Desc: "add [name] [age]: add person",
			Run: func(c *sc.C, args []string) error {
				if len(args) != 2 {
					return sc.UsageError
				}
				name := args[0]
				age, err := strconv.Atoi(args[1])
				if err != nil {
					return err
				}
				return add(name, age)
			},
		},
	}).Run(&sc.C{})
}
