package logger

import "fmt"

func Info(msg string) {
	fmt.Printf("\x1b[37m%s\x1b[0m\n", msg)
}

func Error(msg string) {
	fmt.Printf("\x1b[31m%s\x1b[0m\n", msg)
}