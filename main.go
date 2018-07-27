package main

import (
	"os"
)

func main() {
	cli := &CLI{os.Stdin, os.Stdout, os.Stderr}
	os.Exit(cli.Run(os.Args))
}
