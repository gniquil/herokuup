package main

import (
	"github.com/gniquil/herokuup/herokuup"
	"os"
)

func main() {
	herokuup.NewRunner(os.Args[1]).Run()
}
