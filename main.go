package main

import (
	"github.com/jjzcru/hog/internal/command"
	"os"
)

func main() {
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
