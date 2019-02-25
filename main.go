package main

import (
	"fmt"
)

// Main
func main() {
	a := App{}
	a.Initialise()
	fmt.Println("Running...")
	a.Run(":8000")
}
