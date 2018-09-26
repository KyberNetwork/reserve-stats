package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "hello-world"
	app.Usage = "Hello World"
	app.Action = run
	app.Version = "0.0.1"

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	log.Println("Hello, World!")
	return nil
}
