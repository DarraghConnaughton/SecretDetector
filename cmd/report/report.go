package report

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/types"
	"strconv"
)

func safeDiplayReport(report types.Report, ctx types.Context) {
	log.Println(fmt.Sprintf("[/]start time: %s", strconv.FormatInt(ctx.StartTime, 10)))
	log.Println(fmt.Sprintf("[/]end time: %s [/]total files processed: %d;  in %s time\n# of Potential Secrets Found: %d",
		strconv.FormatInt(ctx.EndTime, 10),
		len(ctx.FilePaths),
		helper.TimeDiff(ctx.StartTime, ctx.EndTime),
		len(report.Secrets)))
}

func writeReport(report types.Report) error {

	dir := filepath.Dir("./report/")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	file, err := os.Create("./report/report.json")
	if err != nil {
		return err
	}
	defer file.Close()
	jsonData, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func Handler(report types.Report, ctx types.Context) error {
	safeDiplayReport(report, ctx)
	return writeReport(report)
}
