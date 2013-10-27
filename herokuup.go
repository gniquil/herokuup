package main

import (
	"github.com/gniquil/herokuup/herokuup"
	"os"
)

func main() {
	herokuup.Run(os.Args[1])
}
