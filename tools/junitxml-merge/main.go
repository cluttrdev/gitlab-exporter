package main

import (
	"context"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter/internal/junitxml"
)

func main() {
	if err := exec(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func exec(ctx context.Context) error {
	cmd := configure()

	args := os.Args[1:]
	opts := []cli.ParseOption{}

	if err := cmd.Parse(args, opts...); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		} else {
			return fmt.Errorf("error parsing arguments: %w", err)
		}
	}

	return cmd.Run(ctx)
}

func configure() *cli.Command {
	fs := flag.NewFlagSet("junitxml-merge", flag.ExitOnError)

	output := fs.String("o", "", "Write to file instead of stdout")

	return &cli.Command{
		Name:       "junitxml-merge",
		ShortHelp:  "Merge junit xml files.",
		ShortUsage: "junitxml-merge [-o <file>] [FILE]...",
		Flags:      fs,
		Exec: func(ctx context.Context, args []string) error {
			filepaths := args

			reports := make([]junitxml.TestReport, 0, len(filepaths))
			for _, filepath := range filepaths {
				file, err := os.Open(filepath)
				if err != nil {
					slog.Error("failed to open file", "error", err, "file", filepath)
					continue
				}

				report, err := junitxml.Parse(file)
				if err != nil {
					slog.Error("failed to parse file", "error", err, "file", filepath)
					continue
				}

				reports = append(reports, report)
			}

			report := junitxml.Merge(reports)

			var ofile *os.File = os.Stdout
			if *output != "" {
				var err error
				ofile, err = os.Create(*output)
				if err != nil {
					return fmt.Errorf("create output file: %w", err)
				}
			}

			encoder := xml.NewEncoder(ofile)
			encoder.Indent("", "\t")
			if err := encoder.Encode(report); err != nil {
				return fmt.Errorf("write file: %w", err)
			}

			return nil
		},
	}
}
