package secretdetection

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/types"
	"sync"
)

// Producer reads lines from multiple files and sends them to the channel
func producer(filePaths []string, ch chan types.Line, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, filePath := range filePaths {
		lineNo := 0
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {

			line := scanner.Text()
			ch <- types.Line{
				Number:   lineNo,
				Contents: line,
				Filename: filePath,
			}
			lineNo += 1
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("[-]Error reading file:", err)
		}
		if err := file.Close(); err != nil {
			fmt.Println("[-]Error closing file:", err)
		}
	}
}

// Consumer reads lines from the channel and prints them along with line numbers
func consumer(ch chan types.Line, outputch chan types.Line, secretPatterns []*regexp.Regexp, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range ch {
		for _, rePattern := range secretPatterns {
			p := helper.DetectPattern(rePattern, line.Contents)
			if len(p) > 0 {
				outputch <- line
			}
		}
	}
}

func amalgamator(ch chan types.Line, outputChan chan types.Report, wg *sync.WaitGroup) {
	defer wg.Done()
	tmp := []types.Line{}
	for line := range ch {
		tmp = append(tmp, line)
	}
	outputChan <- types.Report{
		Timestamp: helper.CurrentTime(),
		Secrets:   tmp,
	}
	close(outputChan)
}

func DetectSecrets(ctx types.Context) (types.Report, error) {
	var report types.Report
	// Channel between producer and consumer
	lineChannel := make(chan types.Line, 50)
	// Channel between consumer and amalgamation goroutine
	outputChannel := make(chan types.Line, 10)
	// Channel between amalgamation goroutine and main goroutine.
	resultChannel := make(chan types.Report, 1)

	var producersWg sync.WaitGroup
	var consumersWg sync.WaitGroup
	var almalWg sync.WaitGroup

	NumOfCPUs = int(math.Min(float64(len(ctx.FilePaths)), float64(NumOfCPUs)))
	splitFiles := helper.SplitFiles(ctx.FilePaths, NumOfCPUs)
	if len(splitFiles) > 0 {
		for i := 0; i < NumOfCPUs; i++ {
			producersWg.Add(1)
			go producer(splitFiles[i], lineChannel, &producersWg)
		}

		for i := 0; i < NumOfCPUs; i++ {
			consumersWg.Add(1)
			go consumer(lineChannel, outputChannel, ctx.SecretPatterns, &consumersWg)
		}

		almalWg.Add(1)
		go amalgamator(outputChannel, resultChannel, &almalWg)

		go func() {
			producersWg.Wait()
			close(lineChannel)
		}()

		go func() {
			consumersWg.Wait()
			close(outputChannel)
		}()
		// Wait for producers and consumers to finish before starting Amalgamator
		producersWg.Wait()
		consumersWg.Wait()
		almalWg.Wait()

		report = <-resultChannel
	}
	return report, nil
}
