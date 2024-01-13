package sequential

import (
	"bufio"
	"log"
	"os"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/types"
)

type SequentialSecretDetector struct {
	types.SecretDetectorInterface
	Report types.Report
	Ctx    types.Context
}

func (ssd *SequentialSecretDetector) StartScan() {
	log.Println("[+] Starting Sequential Secret Detection.")
	var secrets []types.Line
	startTime := helper.CurrentTime()
	for _, filePath := range ssd.Ctx.FilePaths {
		lineNo := 0
		if file, err := os.Open(filePath); err == nil {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				for _, rePattern := range ssd.Ctx.SecretPatterns {
					p := helper.DetectPattern(rePattern, line)
					if len(p) > 0 {
						secrets = append(secrets, types.Line{
							Number:   lineNo,
							Contents: line,
							Filename: filePath,
						})
					}
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
	ssd.Report = types.Report{
		Timestamp: helper.CurrentTime(),
		Secrets:   secrets,
		StartTime: startTime,
		EndTime:   helper.CurrentTime(),
	}
}

func New(ctx types.Context) SequentialSecretDetector {
	return SequentialSecretDetector{
		Ctx:    ctx,
		Report: types.Report{},
	}
}
