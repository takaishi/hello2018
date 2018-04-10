package main

import (
	"github.com/hashicorp/logutils"
	"github.com/mattn/sc"
	"github.com/takaishi/hello2018/grpc_with_transport_credentials/client"
	"github.com/takaishi/hello2018/grpc_with_transport_credentials/server"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
	"fmt"
)

func logLevel() string {
	envLevel := os.Getenv("LOG_LEVEL")
	if envLevel == "" {
		return "WARN"
	} else {
		return envLevel
	}
}

func logOutput() *logutils.LevelFilter {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(logLevel()),
		Writer:   os.Stderr,
	}

	return filter
}

func main() {
	log.SetOutput(logOutput())

	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name: "server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "public-key, p",
					Value: fmt.Sprintf("%s/.ssh/id_rsa,pub", os.Getenv("HOME")),
					Usage: "Set the ssh public key.",
				},
			},
			Action: func(c *cli.Context) error {
				return server.Start(c)
			},
		},
		{
			Name: "client",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "identity-file, i",
					Value: fmt.Sprintf("%s/.ssh/id_rsa", os.Getenv("HOME")),
					Usage: "Set the ssh private key.",
				},
			},
			Subcommands: []cli.Command{
				{
					Name: "add",
					Action: func(c *cli.Context) error {
						if len(c.Args()) != 2 {
							return sc.UsageError
						}

						name := c.Args()[0]
						age, err := strconv.Atoi(c.Args()[1])

						if err != nil {
							return err
						}

						return client.Add(c, name, age)
					},
				},
				{
					Name: "list",
					Action: func(c *cli.Context) error {

						return client.List(c)
					},
				},
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
