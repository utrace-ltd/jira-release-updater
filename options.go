package main

import (
	"github.com/voxelbrain/goptions"
)

type Options struct {
	ComponentName string        `goptions:"-c, --component-name, obligatory, description='Component name'"`
	Version       string        `goptions:"-v, --version, obligatory, description='Version from get change log'"`
	Help          goptions.Help `goptions:"-h, --help, description='Show this help'"`
}
