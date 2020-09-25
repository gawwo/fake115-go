package main

import (
	"github.com/gawwo/fake115-go/executor"
	"testing"
)

func TestScanDirWithOffset(t *testing.T) {
	executor.ScanDirWithOffset("1932902800904947822", 0)
}
