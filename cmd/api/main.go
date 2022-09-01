package main

import (
	"fmt"
)

const version = "1.0"

func main() {
	fmt.Printf("###############################################\n")
	fmt.Printf("# Test assignment RestAPI Server. Version %s #\n", version)
	fmt.Printf("###############################################\n\n")

	app := newApp()

	//Starting web server
	app.serve()
}
