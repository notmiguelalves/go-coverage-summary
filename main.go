package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/cover"
)

type PackageCoverage struct {
	name              string
	totalStatements   int
	coveredStatements int
}

type Coverage map[string]*PackageCoverage

func main() {
	var coverageFile string
	flag.StringVar(&coverageFile, "file", "", "path to coverage file")
	flag.Parse()

	if strings.TrimSpace(coverageFile) == "" {
		log.Println("path to coverage file should not be empty")
		os.Exit(1)
	}

	profiles, err := cover.ParseProfiles(coverageFile)
	if err != nil {
		log.Printf("failed to parse coverage file: %s\n", err.Error())
		os.Exit(1)
	}

	coverage := Coverage{}

	for _, p := range profiles {
		if p == nil {
			log.Println("got nil profile")
			os.Exit(1)
		}

		if !strings.EqualFold(p.Mode, "atomic") {
			log.Println("only coverage profiles in 'atomic' mode are supported")
			os.Exit(1)
		}

		pkg := filepath.Dir(p.FileName)
		if _, ok := coverage[pkg]; !ok {
			coverage[pkg] = &PackageCoverage{
				name:              pkg,
				totalStatements:   0,
				coveredStatements: 0,
			}
		}

		cvg := coverage[pkg]

		for _, b := range p.Blocks {
			cvg.totalStatements += b.NumStmt
			if b.Count > 0 {
				cvg.coveredStatements += b.NumStmt
			}
		}
	}

	writeGithubStepSummary(coverage)
}

func writeGithubStepSummary(coverage Coverage) {
	stepSummaryOutput := os.Getenv("GITHUB_STEP_SUMMARY")
	if stepSummaryOutput == "" {
		log.Println("$GITHUB_STEP_SUMMARY cannot be empty")
		os.Exit(1)
	}

	f, err := os.OpenFile(stepSummaryOutput, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("failed to open github step summary output: %s\n", err.Error())
		os.Exit(1)
	}

	writeLine(f, "\n| Package | Coverage |\n")
	writeLine(f, "| ----- | ----- |\n")
	for pkg, cov := range coverage {
		covPercent := 100 * float64(cov.coveredStatements) / float64(cov.totalStatements)
		writeLine(f, fmt.Sprintf("| `%s` | **%.1f%%** |\n", pkg, covPercent))
	}
}

func writeLine(f *os.File, line string) {
	_, err := f.Write([]byte(line))
	if err != nil {
		log.Printf("failed to write to step output: %s\n", err.Error())
		os.Exit(1)
	}
}
