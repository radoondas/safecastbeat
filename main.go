package main

import (
	"os"

	"github.com/radoondas/safecastbeat/cmd"
	_ "github.com/radoondas/safecastbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
