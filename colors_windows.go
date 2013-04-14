package main

//import "fmt"

const (
	blue  = 0x1
	green = 0x2
	red   = 0x4
)

func Green(input string) string {
	//set := green&0x0F|0&0x0F
	//return fmt.Sprintf("%x%s", set, input)
	return input
}

func Blue(input string) string {
	//set := blue&0x0F | 0&0x0f
	//return fmt.Sprintf("%x%s", set, input)
	return input
}

func Reset() string {
	//reset := (red|blue|green)&0x0F | 0&0x0F
	//return fmt.Sprintf("%x", reset)
	return ""
}
