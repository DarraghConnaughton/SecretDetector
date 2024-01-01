package main

import (
	"fmt"
	"log"
	"os"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/secretdetection"
	"secretdetecion/cmd/types"
)

// defining main function
func main() {

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var context types.Context
	context, err := helper.RetrieveContext()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	startTime := helper.CurrentTime()
	log.Println(fmt.Sprintf("[/]start time: %s", startTime))

	report, err := secretdetection.DetectSecrets(context)
	helper.CheckError(err)

	endTime := helper.CurrentTime()
	log.Println(fmt.Sprintf("[/]end time: %s", endTime))
	log.Println("[/]...")
	log.Println("[/]...")
	log.Println(fmt.Sprintf("[/]total files processed: %d;  in %s time", len(context.FilePaths), helper.TimeDiff(startTime, endTime)))
	log.Println("Amalgamated Report:", report)
	log.Println("# of Potential Secrets Found:", len(report.Secrets))

}
