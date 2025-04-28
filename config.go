package main

import (
	"glstatus/components"
)

type Module struct {
	Producer func(string) string // Component function
	Format   string              // printf-style format
	Argument string              // Component argument
}

var modules = []Module{
	{Producer: components.RamPerc, Format: "Mem: %s ", Argument: ""},
	{Producer: components.RamFree, Format: "Free: %s ", Argument: ""},
}
