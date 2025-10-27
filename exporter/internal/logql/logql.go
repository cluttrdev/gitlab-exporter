package logql

import (
	"bufio"
	"bytes"
)

type MetricQuery struct {
	Name       string
	LineFilter LineFilter
	LabelAdd   map[string]string
}

func Count(log *bytes.Reader, filters []LineFilter) ([]int, error) {
	counts := make([]int, len(filters))

	scanner := bufio.NewScanner(log)
	for scanner.Scan() {
		line := scanner.Bytes()

		for i, filter := range filters {
			if filter.Match(line) {
				counts[i] = counts[i] + 1
			}
		}
	}

	return counts, nil
}
