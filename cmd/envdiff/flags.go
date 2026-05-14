package main

import (
	"flag"
	"fmt"
	"os"
)

// config holds all parsed CLI flags and positional arguments.
type config struct {
	fileA          string
	fileB          string
	format         string
	showValues     bool
	prefix         string
	onlyMissing    bool
	onlyMismatched bool
}

// parseFlags parses os.Args and returns a populated config.
// It writes usage to stderr and exits on error.
func parseFlags() config {
	fs := flag.NewFlagSet("envdiff", flag.ExitOnError)

	format := fs.String("format", "text", "Output format: text or json")
	showValues := fs.Bool("show-values", false, "Show actual values in output")
	prefix := fs.String("prefix", "", "Only compare keys with this prefix")
	onlyMissing := fs.Bool("only-missing", false, "Only show missing keys")
	onlyMismatched := fs.Bool("only-mismatched", false, "Only show mismatched values")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [options] <file-a> <file-b>\n\nOptions:\n")
		fs.PrintDefaults()
	}

	_ = fs.Parse(os.Args[1:])

	args := fs.Args()
	if len(args) != 2 {
		fs.Usage()
		os.Exit(1)
	}

	return config{
		fileA:          args[0],
		fileB:          args[1],
		format:         *format,
		showValues:     *showValues,
		prefix:         *prefix,
		onlyMissing:    *onlyMissing,
		onlyMismatched: *onlyMismatched,
	}
}
