package main

import (
	"log"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/secretdetection/concurrent"
	"secretdetecion/cmd/secretdetection/sequential"
	"secretdetecion/cmd/types"
)

var (
	context    types.Context
	configPath string
	filePath   string
	reportPath string
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	helper.RetrieveFlags(&filePath, &configPath, &reportPath)
	helper.CheckError(
		helper.RetrieveContext(&context, filePath, configPath))

	sequentialDetector := sequential.New(context)
	sequentialDetector.StartScan()
	helper.DiplayReport(sequentialDetector.Report, sequentialDetector.Ctx)

	concurrentDetector := concurrent.New(context)
	concurrentDetector.StartScan()

	helper.CheckError(
		helper.HandleReport(reportPath, concurrentDetector.Report, concurrentDetector.Ctx))
}
