package main

import (
	"log"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/secretdetection"
	"secretdetecion/cmd/types"
)

var (
	context    types.Context
	report     types.Report
	configPath string
	filePath   string
	reportPath string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// Load configuration parameters and secret detection context.
	helper.RetrieveFlags(&filePath, &configPath, &reportPath)
	helper.CheckError(
		helper.RetrieveContext(&context, filePath, configPath))

	// Detect secrets and generate report.
	helper.CheckError(
		secretdetection.DetectSecrets(&report, context))

	// Display and write report.
	helper.CheckError(
		helper.HandleReport(reportPath, report, context))
}
