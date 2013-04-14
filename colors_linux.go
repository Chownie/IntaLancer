package main

import "fmt"

const (
	green = "\x1b[32m"
	blue  = "\x1b[34m"
)

func Green(input string) string {
	return fmt.Sprintf("%s%s", green, input)
}

func Blue(input string) string {
	return fmt.Sprintf("%s%s", blue, input)
}

func Reset() string {
	return "\x1b[0m"
}
