package main

import "fmt"

func main() {
	re := NewRegexMachine("aboba%.")
	fmt.Println(re.Match("aboba."))
}
