package trailerparser

import (
	"bytes"
)

type Trailer struct {
	Key   []byte
	Value []byte
}

// parseCommitTrailer implements git trailer parsing like Gitaly does internally.
//
// Source: gitlab.com/gitlab-org/gitaly@06464fa5/internal/git/trailerparser/trailerparser.go
func Parse(input []byte) []Trailer {
	const (
		maxTrailers  = 64
		maxKeySize   = 128
		maxValueSize = 512
	)

	var trailers []Trailer
	startPos := 0
	max := len(input)

	for startPos < max && len(trailers) < maxTrailers {
		endPos := bytes.IndexByte(input[startPos:], byte('\000'))

		// endPos starts as an index relative to the start position, so we need
		// to convert this to an absolute index before we use it for slicing.
		if endPos >= 0 {
			endPos = startPos + endPos
		} else {
			endPos = max
		}

		sepPos := bytes.IndexByte(input[startPos:endPos], byte(':'))

		if sepPos > 0 {
			key := input[startPos : startPos+sepPos]
			value := bytes.TrimSpace(input[startPos+sepPos+1 : endPos])

			if len(key) > maxKeySize || len(value) > maxValueSize {
				break
			}

			trailers = append(trailers, Trailer{
				Key:   key,
				Value: value,
			})
		}

		startPos = endPos + 1
	}

	return trailers
}

func ExtractPotentialTrailerLines(message []byte) [][]byte {
	blocks := bytes.Split(message, []byte("\n\n")) // split on blank lines
	if len(blocks) < 2 {                           // expecting at least title + trailer lines
		return nil
	}

	// split last message block lines
	return bytes.Split(blocks[len(blocks)-1], []byte("\n"))
}
