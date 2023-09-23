package logger

import (
	"fmt"
	"log"
)

func LogBlue(format string, a ...interface{}) {
	fmt.Printf("\033[94m"+format+"\033[0m", a...)
	fmt.Println()
}

func LogGreen(format string, a ...interface{}) {
	fmt.Printf("\033[32m"+format+"\033[0m", a...)
	fmt.Println()
}

func LogFatalRedAndExit(format string, a ...interface{}) {
	log.Fatalf("\033[91m"+format+"\033[0m", a)
}

func LogError(format string, a ...interface{}) {
	fmt.Printf("\033[91m"+format+"\033[0m", a...)
	fmt.Println()
}
