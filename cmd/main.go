package main

import (
	"fmt"
	"log"
	"os"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/secretdetection"
	"secretdetecion/cmd/types"
	"strconv"
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
	log.Println(fmt.Sprintf("[/]start time: %s", strconv.FormatInt(startTime, 10)))
	report, err := secretdetection.DetectSecrets(context)
	helper.CheckError(err)

	endTime := helper.CurrentTime()
	log.Println(fmt.Sprintf("[/]end time: %s [/]total files processed: %d;  in %s time\n# of Potential Secrets Found: %d",
		strconv.FormatInt(endTime, 10),
		len(context.FilePaths),
		helper.TimeDiff(startTime, endTime),
		len(report.Secrets)))
}
