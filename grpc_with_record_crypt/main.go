package main

import (
	"github.com/hashicorp/logutils"
	"github.com/takaishi/hello2018/grpc_with_record_crypt/command"
	"github.com/urfave/cli"
	"log"
	"os"
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
		cli.BoolFlag{
			Name: "secure",
		},
	}
	app.Commands = []cli.Command{
		{
			Name: "server",
			Action: func(c *cli.Context) error {
				return command.StartServer(c, c.Bool("secure"))
			},
		},
		{
			Name: "client",
			Action: func(c *cli.Context) error {
				return command.StartClient(c, c.Bool("secure"))
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
