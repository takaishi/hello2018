package main

import (
	"fmt"
	"io"
	"strconv"

	"github.com/mattn/sc"
	pb "github.com/takaishi/hello2018/grpc_with_stns/protocol"
	sshTC2 "github.com/takaishi/hello2018/grpc_with_stns/sshTC"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/hashicorp/logutils"
	"os"
	"log"
)

func add(name string, age int) error {
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

func list() error {
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
		fmt.Println(err)
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

func logLevel() string {
	envLevel := os.Getenv("LOG_LEVEL")
	if envLevel == "" {
		return "WARN"
	} else {
		return envLevel
	}
}

func logOutput() (*logutils.LevelFilter){
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(logLevel()),
		Writer: os.Stderr,
	}

	return filter
}

func main() {
	log.SetOutput(logOutput())

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
