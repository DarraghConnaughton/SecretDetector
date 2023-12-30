package secretdetection

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"secretdetecion/cmd/helper"
	"sync"
)

// Producer reads lines from multiple files and sends them to the channel
func Producer(filePaths []string, ch chan Line, wg *sync.WaitGroup) {
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
			ch <- Line{
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
func Consumer(ch chan Line, outputch chan Line, secretPatterns []*regexp.Regexp, wg *sync.WaitGroup) {
	defer wg.Done()
	for line := range ch {
		for _, rePattern := range secretPatterns {
			p := DetectPattern(rePattern, line.Contents)
			if len(p) > 0 {
				outputch <- line
			}
		}
	}
}

func Amalgamate(ch chan Line, outputChan chan Report, wg *sync.WaitGroup) {
	defer wg.Done()
	tmp := []Line{}
	for line := range ch {
		tmp = append(tmp, line)
	}

	outputChan <- Report{
		Timestamp: helper.CurrentTime(),
		Secrets:   tmp,
	}
	close(outputChan)
}
