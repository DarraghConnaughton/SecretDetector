package main

import (
	"log"
	"secretdetecion/cmd/helper"
	r "secretdetecion/cmd/report"
	"secretdetecion/cmd/secretdetection"
	"secretdetecion/cmd/types"
)

// defining main function
func main() {
	var configPath, filePath string
	var context types.Context
	var err error

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	helper.RetrieveFlags(&filePath, &configPath)

	context, err = helper.RetrieveContext(filePath, configPath)
	helper.CheckError(err)

	context.StartTime = helper.CurrentTime()
	report, err := secretdetection.DetectSecrets(context)
	helper.CheckError(err)

	context.EndTime = helper.CurrentTime()
	err = r.Handler(report, context)
	helper.CheckError(err)
}
