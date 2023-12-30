package secretdetection

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"log"
	"os"
	"regexp"
	"secretdetecion/cmd/helper"
)

func DetectPattern(regex *regexp.Regexp, line string) []string {
	return regex.FindAllString(line, -1)
}

func GetToml(filename string) []*regexp.Regexp {
	log.Println(fmt.Sprintf("[+]retrieving secret patterns from %s", filename))
	info, err := os.Stat(filename)
	helper.CheckError(err)

	log.Println(fmt.Sprintf("[+]file detected. Metadata: %s", info))

	dat, err := os.ReadFile(filename)
	helper.CheckError(err)

	var cfg TomlConfig
	err = toml.Unmarshal(dat, &cfg)

	compiledRegex := []*regexp.Regexp{}
	for _, v := range cfg.Rules {
		compiledRegex = append(compiledRegex, regexp.MustCompile(v.Regex))
	}
	return compiledRegex
}
