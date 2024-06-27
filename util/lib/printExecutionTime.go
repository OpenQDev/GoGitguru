package lib

import (
	"fmt"
	"time"
)

func PrintExecutionTime(startTime int64, functionName string, repoUrl string) {
	fmt.Println("Execution time for", functionName, ":", time.Now().Unix()-startTime, "for", repoUrl)
}
