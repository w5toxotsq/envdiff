package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/report"
)

func main() {
	var (
		format     = flag.String("format", "text", "Output format: text or json")
		showValues = flag.Bool("show-values", false, "Show actual values in output (may expose secrets)")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [options] <file1> <file2>\n\n")
		fmt.Fprintf(os.Stderr, "Compare two .env files and report differences.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	leftPath, rightPath := args[0], args[1]

	left, err := parser.ParseFile(leftPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", leftPath, err)
		os.Exit(1)
	}

	right, err := parser.ParseFile(rightPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", rightPath, err)
		os.Exit(1)
	}

	results := diff.Compare(left, right)

	opts := report.DefaultOptions()
	opts.Format = *format
	opts.ShowValues = *showValues
	opts.LeftLabel = leftPath
	opts.RightLabel = rightPath

	output, err := report.Render(results, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error rendering report: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(output)

	if len(results) > 0 {
		os.Exit(2)
	}
}
