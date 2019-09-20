package main

import (
	"github.com/VanMoof/gopenapi/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}
