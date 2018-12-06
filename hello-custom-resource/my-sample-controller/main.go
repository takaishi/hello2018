package main

import (
	"github.com/urfave/cli"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "kubeconfig",
			Value: "",
		},
		cli.StringFlag{
			Name:  "masterURL",
			Value: "",
		},
	}

	app.Action = func(c *cli.Context) error {
		return action(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	cfg, err := clientcmd.BuildConfigFromFlags(c.String("masterURL"), c.String("kubeconfig"))
	if err != nil {
		return err
	}
	log.Printf("%+v\n", cfg)
	return nil
}
