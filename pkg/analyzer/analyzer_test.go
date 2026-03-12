package analyzer_test

import (
	"github.com/LeksinMaksim/loglinter/pkg/analyzer"
	"golang.org/x/tools/go/analysis/analysistest"
	"testing"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer.Analyzer, "a")

}
