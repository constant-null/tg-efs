package main

import (
	_ "embed"
)

func main() {
	go runWebApp()
	runTGBot()
}
