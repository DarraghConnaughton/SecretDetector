package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/secretdetection"
	"sync"
)

// defining main function
func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	secretPatterns := secretdetection.GetToml(fmt.Sprintf("%s/cmd/secretdetection/data/secretpatterns.toml", pwd))
	log.Println(fmt.Sprintf("[/]secret patterns loaded: %d", len(secretPatterns)))
	files, err := helper.CollectFiles("/Users/darraghconnaughton/Github/thirdParty/SecretsTest")
	helper.CheckError(err)

	startTime := helper.CurrentTime()
	log.Println(fmt.Sprintf("[/]start time: %s", startTime))

	numCPU := runtime.NumCPU() / 2

	var report secretdetection.Report
	// Channel between producer and consumer
	lineChannel := make(chan secretdetection.Line, 50)
	// Channel between consumer and amalgamation goroutine
	outputChannel := make(chan secretdetection.Line, 10)
	// Channel between amalgamation goroutine and main goroutine.
	resultChannel := make(chan secretdetection.Report, 1)

	var producersWg sync.WaitGroup
	var consumersWg sync.WaitGroup
	var almalWg sync.WaitGroup

	splitFiles := helper.SplitFiles(files, numCPU)
	for i := 0; i < numCPU; i++ {
		producersWg.Add(1)
		go secretdetection.Producer(splitFiles[i], lineChannel, &producersWg)
	}

	for i := 0; i < numCPU; i++ {
		consumersWg.Add(1)
		go secretdetection.Consumer(lineChannel, outputChannel, secretPatterns, &consumersWg)
	}

	almalWg.Add(1)
	go secretdetection.Amalgamate(outputChannel, resultChannel, &almalWg)

	go func() {
		producersWg.Wait()
		close(lineChannel)
	}()

	go func() {
		consumersWg.Wait()
		close(outputChannel)
	}()
	// Wait for producers and consumers to finish before starting Amalgamate
	producersWg.Wait()
	consumersWg.Wait()
	almalWg.Wait()

	report = <-resultChannel
	endTime := helper.CurrentTime()
	log.Println(fmt.Sprintf("[/]end time: %s", endTime))
	log.Println("[/]...")
	log.Println("[/]...")
	log.Println(fmt.Sprintf("[/]total files processed: %d;  in %s time", len(files), helper.TimeDiff(startTime, endTime)))
	log.Println("Amalgamated Report:", report)
	log.Println("# of Potential Secrets Found:", len(report.Secrets))

}
