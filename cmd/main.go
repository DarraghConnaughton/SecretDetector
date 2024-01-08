package main

import (
	"log"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/secretdetection"
	"secretdetecion/cmd/types"
)

// defining main function
func main() {
	var configPath, filePath, reportPath string
	var context types.Context
	var err error

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	helper.RetrieveFlags(&filePath, &configPath, &reportPath)

	context, err = helper.RetrieveContext(filePath, configPath)
	helper.CheckError(err)

	context.StartTime = helper.CurrentTime()
	report, err := secretdetection.DetectSecrets(context)
	helper.CheckError(err)

	context.EndTime = helper.CurrentTime()
	err = helper.HandleReport(reportPath, report, context)
	helper.CheckError(err)
}
