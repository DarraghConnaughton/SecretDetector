package helper

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"secretdetecion/cmd/types"
	"strings"
	"time"
)

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func collectFiles(startDir string) []string {
	var files []string
	filepath.Walk(startDir, func(fp string, fi os.FileInfo, err error) error {
		if !fi.IsDir() && !strings.Contains(fp, "vendor/") {
			files = append(files, fp)
		}
		return nil
	})
	return files
}

func CurrentTime() int64 {
	return time.Now().Unix()
}

func TimeDiff(t1 int64, t2 int64) time.Duration {
	return time.Unix(t2, 0).Sub(time.Unix(t1, 0))
}

func SplitFiles(files []string, num int) [][]string {
	if num <= 0 {
		num = 1
	}
	totalFiles := len(files)
	if num <= 0 || totalFiles == 0 {
		return nil
	}
	partSize := int(math.Ceil(float64(totalFiles) / float64(num)))
	result := make([][]string, 0, num)
	for i := 0; i < totalFiles; i += partSize {
		end := i + partSize
		if end > totalFiles {
			end = totalFiles
		}
		result = append(result, files[i:end])
	}
	return result
}

func RetrieveFlags(filepath *string, configpath *string, reportPath *string) {
	flag.StringVar(filepath, "filepath", "./", "start path for recursive search. \n[default] current directory.")
	flag.StringVar(configpath, "configpath", "/cmd/data/secretpatterns.toml", "regex for known secret patterns.")
	flag.StringVar(reportPath, "reportpath", "./report.json", "where to write report.")
	flag.Parse()
}

func RetrieveContext(tContext *types.Context, searchStartDir string, configPath string) error {
	var err error
	tContext.SecretPatterns, err = GetToml(configPath)
	if err != nil {
		return err
	}
	// Retrieving files paths is the more expensive operation. No point initiating this process if the
	// config is invalid.
	tContext.FilePaths = collectFiles(searchStartDir)

	log.Println(fmt.Sprintf("[/]secret patterns loaded: %d", len(tContext.SecretPatterns)))
	return nil
}

func DetectPattern(regex *regexp.Regexp, line string) []string {
	return regex.FindAllString(line, -1)
}

func GetToml(filename string) ([]*regexp.Regexp, error) {
	log.Println(fmt.Sprintf("[+]retrieving secret patterns from %s", filename))
	_, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	dat, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg types.TomlConfig
	err = toml.Unmarshal(dat, &cfg)
	if err != nil {
		return nil, err
	}

	compiledRegex := []*regexp.Regexp{}
	for _, v := range cfg.Rules {
		compiledRegex = append(compiledRegex, regexp.MustCompile(v.Regex))
	}
	return compiledRegex, nil
}

func DiplayReport(report types.Report, ctx types.Context) {
	log.Println(fmt.Sprintf("[/]start time: %d", report.StartTime))
	log.Println(fmt.Sprintf("[/]end time: %d [/]total files processed: %d;  in %s time\n# of Potential Secrets Found: %d",
		report.EndTime,
		len(ctx.FilePaths),
		TimeDiff(report.StartTime, report.EndTime),
		len(report.Secrets)))
}

func writeReport(writePath string, report types.Report) error {
	file, err := os.Create(writePath)
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

func HandleReport(writePath string, report types.Report, ctx types.Context) error {
	DiplayReport(report, ctx)
	return writeReport(writePath, report)
}
