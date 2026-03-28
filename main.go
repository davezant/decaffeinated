/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"decaffeinated/internal/dprocesses"
)

func main() {
	m := dprocesses.NewMonitor()
	m.IncludeCurrentProcesses()
	}
