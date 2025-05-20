package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Prompt for name
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("What's your name? ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	// Greet and show length
	fmt.Printf("Hello, %s!\n", name)
	fmt.Printf("Your name has %d characters.\n", len(name))
}
