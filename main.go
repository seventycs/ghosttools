package main

import (
	"ghosttools/cmd"
	"ghosttools/internal/ui"
	"ghosttools/internal/updater"
)

func main() {
	go updater.Check()
	ui.Startup()
	ui.AnimatedBanner()
	cmd.Run()
}