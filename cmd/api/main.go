package main

import (
	"fmt"
)

const version = "0.1"

func main() {
	fmt.Printf("###############################################\n")
	fmt.Printf("# Test assignment RestAPI Server. Version %s #\n", version)
	fmt.Printf("###############################################\n\n")

	app := newApp()

	//Starting web server
	app.serve()
}
