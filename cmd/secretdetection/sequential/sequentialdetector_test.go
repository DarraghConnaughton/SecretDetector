package sequential

import (
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/types"
	"testing"
)

func TestSequentialSecretDetection(t *testing.T) {
	patterns, err := helper.GetToml("./testdata/secretpatterns.toml")
	if err != nil {
		t.Errorf("sequentialdetector: unable to load test patterns.")
	}
	sequentialDetector := New(types.Context{
		FilePaths:      []string{"./testdata/testsecrets/forgottenAboutCredentials.txt"},
		SecretPatterns: patterns,
	})
	sequentialDetector.StartScan()
	if len(sequentialDetector.Report.Secrets) != 1 {
		t.Errorf("sequentialdetector: expected number of secrets deviate from expected amount.")
	}
}
