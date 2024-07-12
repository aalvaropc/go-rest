package main

import (
	"fmt"
)

type Task struct {
	ID      int    `json:ID`
	Name    string `json:ID`
	Content string `json:ID`
}

func main() {
	fmt.Println("Hello World!")
}
