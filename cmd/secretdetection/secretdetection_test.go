package secretdetection

import (
	"secretdetecion/cmd/helper"
	"secretdetecion/cmd/types"
	"testing"
)

func TestSecretDetection(t *testing.T) {
	report, err := DetectSecrets(types.Context{
		FilePaths:      []string{"./testdata/testsecrets/forgottenAboutCredentials.txt"},
		SecretPatterns: nil,
	})
	if err != nil {
		t.Errorf("")
	}
	if len(report.Secrets) != 0 {
		t.Errorf("")
	}

	patterns, err := helper.GetToml("./testdata/secretpatterns.toml")
	if err != nil {
		t.Errorf("")
	}

	report, err = DetectSecrets(types.Context{
		FilePaths:      []string{"./testdata/testsecrets/forgottenAboutCredentials.txt"},
		SecretPatterns: patterns,
	})
	if err != nil {
		t.Errorf("")
	}
	if len(report.Secrets) != 1 {
		t.Errorf("")
	}
}
