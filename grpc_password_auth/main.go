package main

import (
	"github.com/hashicorp/logutils"
	"github.com/mattn/sc"
	"github.com/takaishi/hello2018/grpc_password_auth/client"
	"github.com/takaishi/hello2018/grpc_password_auth/server"
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
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
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "username, u",
			Value: "alice",
			Usage: "Set the username to authenticate.",
		},
		cli.StringFlag{
			Name:  "password, p",
			Value: "password",
			Usage: "Set the password to authenticate.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "server",
			Action: func(c *cli.Context) error {
				return server.Start(c)
			},
		},
		{
			Name: "client",
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
