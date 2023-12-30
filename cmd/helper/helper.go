package helper

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"
)

func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

func CollectFiles(startDir string) ([]string, error) {
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
