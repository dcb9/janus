package main

import (
	"log"

	"github.com/dcb9/janus/cli"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	cli.Run()
}
