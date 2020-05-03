package main

import (
	OS "os"

	v "github.com/jjzcru/hog/internal/command/version"
	"github.com/jjzcru/hog/internal/command"
)

var version = ""
var os = ""
var arch = ""
var commit = ""
var date = ""
var goversion = ""

func main() {
	v.SetVersion(version, os, arch, commit, date, goversion)
	err := command.Execute()
	if err != nil {
		OS.Exit(1)
	}
}
