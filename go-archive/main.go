package main

import (
	"fmt"
	"github.com/mholt/archiver"
)

func main() {
	fmt.Println("AAAAAAAAAAAA")
	archiver.TarGz.Make("/tmp/helloworld.tar.gz", []string{"helloworld.dig"})

}
