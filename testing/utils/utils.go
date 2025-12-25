package utils

import (
	"os"
	"strings"
)

func TestsRootDir() string {
	wd, _ := os.Getwd()

	return strings.SplitAfter(wd, "testing")[0]
}
