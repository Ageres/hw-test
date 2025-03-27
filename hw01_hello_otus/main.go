package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	// Place your code here.
	inputString := "Hello, OTUS!"
	reversedString := reverse.String(inputString)
	fmt.Println(reversedString)
}
