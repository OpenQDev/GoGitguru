package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	dbscripts "github.com/OpenQDev/GoGitguru/dbscripts/src"
)

func main() {
	fmt.Println("starting db script")
	dbscripts.ConsolidateFileName()
}

func setupSignalHandler(stopChan chan<- struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		<-sigChan
		close(stopChan)
	}()
}
