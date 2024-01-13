package concurrent

import (
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/types"
	"testing"
)

func TestSecretDetection(t *testing.T) {
	concurrentDetector := New(types.Context{
		FilePaths:      []string{"./testdata/testsecrets/forgottenAboutCredentials.txt"},
		SecretPatterns: nil,
	})
	concurrentDetector.StartScan()
	if len(concurrentDetector.Report.Secrets) != 0 {
		t.Errorf("concurrentdetector: expected no entries to be retrieved during first scan but secrets were found.")
	}

	patterns, err := helper.GetToml("./testdata/secretpatterns.toml")
	if err != nil {
		t.Errorf("concurrentdetector: unable to load test patterns.")
	}

	concurrentDetector = New(types.Context{
		FilePaths:      []string{"./testdata/testsecrets/forgottenAboutCredentials.txt"},
		SecretPatterns: patterns,
	})
	concurrentDetector.StartScan()
	if len(concurrentDetector.Report.Secrets) != 1 {
		t.Errorf("concurrentdetector: expected number of secrets deviate from expected amount.")
	}
}
