package main

import (
	clientset "github.com/takaishi/hello2018/hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	app := cli.NewApp()
	app.Flags = []cli.Flag{}

	app.Action = func(c *cli.Context) error {
		return action(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func action(c *cli.Context) error {
	log.Printf("START!")
	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Printf(err.Error())
	}
	client, err := clientset.NewForConfig(cfg)
	if err != nil {
		log.Printf(err.Error())
	}

	for {
		foos, err := client.SamplecontrollerV1alpha().Foos("default").List(v1.ListOptions{})
		if err != nil {
			log.Printf(err.Error())
			continue
		}

		fmt.Printf("%+v\n", foos)

		time.Sleep(10 * time.Second)
	}
	return nil
}
