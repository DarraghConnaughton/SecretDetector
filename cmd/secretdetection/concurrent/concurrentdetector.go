package concurrent

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"runtime"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/types"
	"sync"
)

var NumOfCPUs = runtime.NumCPU() / 2

type ConcurrentSecretDetector struct {
	types.SecretDetectorInterface
	Report types.Report
	Ctx    types.Context
}

// Producer reads lines from multiple files and sends them to the channel
func producer(filePaths []string, ch chan types.Line, wg *sync.WaitGroup) {
	//fmt.Println("producer coming online!")

	defer wg.Done()
	for _, filePath := range filePaths {
		lineNo := 0
		file, err := os.Open(filePath)
		if err != nil {
			log.Println("Error opening file:", err)
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
			log.Println("[-]Error reading file:", err)
		}
		if err := file.Close(); err != nil {
			log.Println("[-]Error closing file:", err)
		}
	}
}

// Consumer reads lines from the channel and prints them along with line numbers
func consumer(ch chan types.Line, outputch chan types.Line, secretPatterns []*regexp.Regexp, wg *sync.WaitGroup) {
	//fmt.Println("consumer coming online!")

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
	//fmt.Println("amalgamator coming online!")

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

func (csd *ConcurrentSecretDetector) StartScan() {
	log.Println("[+] Starting Concurrent Secret Detection.")

	var starttime int64
	// Channel between producer and consumer
	lineChannel := make(chan types.Line, 50)
	// Channel between consumer and amalgamation goroutine
	outputChannel := make(chan types.Line, 10)
	// Channel between amalgamation goroutine and main goroutine.
	resultChannel := make(chan types.Report, 1)

	var producersWg sync.WaitGroup
	var consumersWg sync.WaitGroup
	var almalWg sync.WaitGroup

	NumOfCPUs = int(math.Min(float64(len(csd.Ctx.FilePaths)), float64(NumOfCPUs)))
	splitFiles := helper.SplitFiles(csd.Ctx.FilePaths, NumOfCPUs)
	if len(splitFiles) > 0 {

		//Benchmarking should start here
		starttime = helper.CurrentTime()
		for i := 0; i < NumOfCPUs; i++ {
			producersWg.Add(1)
			go producer(splitFiles[i], lineChannel, &producersWg)
		}

		for i := 0; i < NumOfCPUs; i++ {
			consumersWg.Add(1)
			go consumer(lineChannel, outputChannel, csd.Ctx.SecretPatterns, &consumersWg)
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

		csd.Report = <-resultChannel
	}
	csd.Report.EndTime = helper.CurrentTime()
	csd.Report.StartTime = starttime
}

func New(ctx types.Context) ConcurrentSecretDetector {
	return ConcurrentSecretDetector{
		Ctx:    ctx,
		Report: types.Report{},
	}
}
