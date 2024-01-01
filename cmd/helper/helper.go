package helper

import (
	"flag"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"secretdetecion/cmd/types"
	"time"
)

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func collectFiles(startDir string) ([]string, error) {
	var files []string

	err := filepath.Walk(startDir, func(fp string, fi os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if !fi.IsDir() {
			files = append(files, fp)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path %v: %v", startDir, err)
	}

	return files, nil
}

func CurrentTime() int64 {
	currentTime := time.Now()
	return currentTime.Unix()
}

func TimeDiff(t1 int64, t2 int64) time.Duration {
	return time.Unix(t2, 0).Sub(time.Unix(t1, 0))
}

func PrintReport() {

}

func SplitFiles(files []string, num int) [][]string {
	totalFiles := len(files)

	if num <= 0 || totalFiles == 0 {
		return nil
	}

	// Calculate the size of each part
	partSize := int(math.Ceil(float64(totalFiles) / float64(num)))

	// Initialize the result slices
	result := make([][]string, 0, num)

	// Split the files into parts
	for i := 0; i < totalFiles; i += partSize {
		end := i + partSize
		if end > totalFiles {
			end = totalFiles
		}
		result = append(result, files[i:end])
	}

	return result
}

func retrieveFilePath() string {
	var tfilePath string
	flag.StringVar(&tfilePath, "filepath", "./", "start path for recursive search.")
	flag.Parse()
	return tfilePath
}

func RetrieveContext() (types.Context, error) {
	var tContext types.Context
	pwd, err := os.Getwd()
	if err != nil {
		return tContext, err
	}
	tContext.SecretPatterns, err = getToml(fmt.Sprintf("%s/cmd/secretdetection/data/secretpatterns.toml", pwd))
	if err != nil {
		return tContext, err
	}

	tContext.FilePaths, err = collectFiles(retrieveFilePath())
	if err != nil {
		return tContext, err
	}

	log.Println(fmt.Sprintf("[/]secret patterns loaded: %d", len(tContext.SecretPatterns)))
	return tContext, nil
}

func DetectPattern(regex *regexp.Regexp, line string) []string {
	return regex.FindAllString(line, -1)
}

func getToml(filename string) ([]*regexp.Regexp, error) {
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
