package analyzer

import "golang.org/x/tools/go/analysis"

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{Analyzer}, nil
}
