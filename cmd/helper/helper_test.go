package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"secretdetecion/cmd/types"
	"testing"
	"time"
)

func readTestReportFile(filePath string) (types.Report, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return types.Report{}, err
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return types.Report{}, err
	}
	var report types.Report
	err = json.Unmarshal(data, &report)
	if err != nil {
		return types.Report{}, err
	}
	err = os.Remove(filePath)
	if err != nil {
		log.Println(fmt.Sprintf("warning: test cleanup failed %s", err.Error()))
	}
	return report, nil
}

func TestSplitFiles(t *testing.T) {
	tmp := SplitFiles([]string{"1", "2", "3"}, 1)
	if len(tmp) != 1 {
		t.Errorf("")
	}
	if len(tmp[0]) != 3 {
		t.Errorf("")
	}

	tmp = SplitFiles([]string{"1", "2", "3"}, 0)
	if len(tmp) != 1 {
		t.Errorf("")
	}
	if len(tmp[0]) != 3 {
		t.Errorf("")
	}

	tmp = SplitFiles([]string{"1", "2", "3"}, 2)
	if len(tmp) != 2 {
		t.Errorf("")
	}
	if len(tmp[0]) != 2 {
		t.Errorf("")
	}
	if len(tmp[1]) != 1 {
		t.Errorf("")
	}
}
func TestTimeFunctions(t *testing.T) {
	t1 := CurrentTime()
	time.Sleep(1 * time.Second)

	if t1 == CurrentTime() {
		t.Errorf("")
	}
	if TimeDiff(t1, CurrentTime()) < 0 {
		t.Errorf("")
	}
}

func TestRetrieveContext(t *testing.T) {
	var ctx types.Context
	err := RetrieveContext(&ctx, "/", "./notreal")
	if err == nil {
		t.Errorf("")
	}
	if len(ctx.SecretPatterns) != 0 {
		t.Errorf("")
	}
	if len(ctx.FilePaths) != 0 {
		t.Errorf("")
	}

	err = RetrieveContext(&ctx, "./", "./testdata/secretpatterns.toml")
	if err != nil {
		t.Errorf("")
	}
	if len(ctx.SecretPatterns) != 121 {
		t.Errorf("")
	}
	if len(ctx.FilePaths) != 3 {
		t.Errorf("")
	}
}

func TestRetrieveFlagsAndCheckError(t *testing.T) {
	var s1, s2, s3 string
	RetrieveFlags(&s1, &s2, &s3)

	if s1 != "./" {
		t.Errorf("fail: expected default secret path.")
	}
	if s2 != "./data/secretpatterns.toml" {
		t.Errorf("fail: expected default secret patterns path.")
	}
	if s3 != "./report.json" {
		t.Errorf("fail: expected default report path.")
	}
	CheckError(nil)
}

func TestCheckError(t *testing.T) {
	done := make(chan string, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- "info: goroutine panicked as expected."
			} else {
				done <- "error: goroutine did not panic"
			}
			close(done)
		}()
		CheckError(errors.New("woops"))
	}()
	result := <-done
	if result == "error: goroutine did not panic" {
		t.Errorf("")
	}
}

func TestHandleReport(t *testing.T) {
	err := HandleReport("/this/is/an/invalid/file/path\n", types.Report{}, types.Context{})
	fmt.Println(err)
	if err == nil {
		t.Errorf("")
	}
	err = HandleReport("./report.json", types.Report{}, types.Context{})
	if err != nil {
		t.Errorf("")
	}

	tReport, err := readTestReportFile("./report.json")
	if err != nil {
		t.Errorf("Error reading report file: %v", err)
	}
	if len(tReport.Secrets) != 0 {
		t.Errorf("expected no secrets to be found.")
	}
	fmt.Println(err)
}

func TestDetectPattern(t *testing.T) {
	if DetectPattern(regexp.MustCompile(`\d+`), "abc 123 def")[0] != "123" {
		t.Errorf("expected no secrets to be found.")
	}
}
