package main

import (
	"fmt"
	"os"

	"github.com/soramon0/jarvis/pkg/commands"
)

func main() {
	if err := commands.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
