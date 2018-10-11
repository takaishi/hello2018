package command

import (
	pb "github.com/takaishi/hello2018/grpc_with_record_crypt/protocol"
	"github.com/urfave/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

func StartClient(c *cli.Context) error {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("127.0.0.1:11111", opts...)

	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	req := &pb.Request{
		Msg: "Hello!!!",
	}

	resp, err := client.Send(context.Background(), req)
	log.Printf("resp = %s\n", resp.Msg)

	return err
}
