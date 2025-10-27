// Copyright 2014 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expfmt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

// A stateFn is a function that represents a state in a state machine. By
// executing it, the state is progressed to the next state. The stateFn returns
// another stateFn, which represents the new state. The end state is represented
// by nil.
type stateFn func() stateFn

// ParseError signals errors while parsing the simple and flat text-based
// exchange format.
type ParseError struct {
	Line int
	Msg  string
}

// Error implements the error interface.
func (e ParseError) Error() string {
	return fmt.Sprintf("text format parsing error in line %d: %s", e.Line, e.Msg)
}

// TextParser is used to parse the simple and flat text-based exchange format. Its
// zero value is ready to use.
type TextParser struct {
	buf          *bufio.Reader // Where the parsed input is read through.
	err          error         // Most recent error.
	currentByte  byte          // The most recent byte read.
	currentToken bytes.Buffer  // Re-used each time a token has to be gathered from multiple bytes.

	metric           *Metric
	currentLabelPair *MetricLabelPair
}

type Metric struct {
	Name        string
	Labels      []MetricLabelPair
	Value       float64
	TimestampMs int64
}

type MetricLabelPair struct {
	Name  string
	Value string
}

// TextToMetricFamilies reads 'in' as the simple and flat text-based exchange
// format and creates MetricFamily proto messages. It returns the MetricFamily
// proto messages in a map where the metric names are the keys, along with any
// error encountered.
func (p *TextParser) LineToMetric(line []byte) (*Metric, error) {
	if len(line) == 0 {
		return nil, nil
	}
	p.reset(bytes.NewReader(line))
	for nextState := p.startOfLine; nextState != nil; nextState = nextState() {
		// Magic happens here...
	}

	// If p.err is io.EOF now, we have run into a premature end of the input
	// stream. Turn this error into something nicer and more
	// meaningful. (io.EOF is often used as a signal for the legitimate end
	// of an input stream.)
	if p.err == io.EOF {
		// p.parseError("unexpected end of input stream")
		p.err = nil
	}
	return p.metric, p.err
}

func (p *TextParser) reset(in io.Reader) {
	p.metric = &Metric{}
	if p.buf == nil {
		p.buf = bufio.NewReader(in)
	} else {
		p.buf.Reset(in)
	}
	p.err = nil
}

// startOfLine represents the state where the next byte read from p.buf is the
// start of a line (or whitespace leading up to it).
func (p *TextParser) startOfLine() stateFn {
	if p.skipBlankTab(); p.err != nil {
		// This is the only place that we expect to see io.EOF,
		// which is not an error but the signal that we are done.
		// Any other error that happens to align with the start of
		// a line is still an error.
		if p.err == io.EOF {
			p.err = nil
		}
		return nil
	}
	return p.readingMetricName
}

// readingMetricName represents the state where the last byte read (now in
// p.currentByte) is the first byte of a metric name.
func (p *TextParser) readingMetricName() stateFn {
	if p.readTokenAsMetricName(); p.err != nil {
		return nil
	}
	if p.currentToken.Len() == 0 {
		p.parseError("invalid metric name")
		return nil
	}
	p.metric.Name = p.currentToken.String()
	// Do not append the newly created currentMetric to
	// currentMF.Metric right now. First wait if this is a summary,
	// and the metric exists already, which we can only know after
	// having read all the labels.
	if p.skipBlankTabIfCurrentBlankTab(); p.err != nil {
		return nil // Unexpected end of input.
	}
	return p.readingLabels
}

// readingLabels represents the state where the last byte read (now in
// p.currentByte) is either the first byte of the label set (i.e. a '{'), or the
// first byte of the value (otherwise).
func (p *TextParser) readingLabels() stateFn {
	if p.currentByte != '{' {
		return p.readingValue
	}
	return p.startLabelName
}

// startLabelName represents the state where the next byte read from p.buf is
// the start of a label name (or whitespace leading up to it).
func (p *TextParser) startLabelName() stateFn {
	if p.skipBlankTab(); p.err != nil {
		return nil // Unexpected end of input.
	}
	if p.currentByte == '}' {
		if p.skipBlankTab(); p.err != nil {
			return nil // Unexpected end of input.
		}
		return p.readingValue
	}
	if p.readTokenAsLabelName(); p.err != nil {
		return nil // Unexpected end of input.
	}
	if p.currentToken.Len() == 0 {
		p.parseError(fmt.Sprintf("invalid label name for metric %q", p.metric.Name))
		return nil
	}
	p.currentLabelPair = &MetricLabelPair{Name: p.currentToken.String()}
	if p.skipBlankTabIfCurrentBlankTab(); p.err != nil {
		return nil // Unexpected end of input.
	}
	if p.currentByte != '=' {
		p.parseError(fmt.Sprintf("expected '=' after label name, found %q", p.currentByte))
		return nil
	}
	// Check for duplicate label names.
	for _, l := range p.metric.Labels {
		if l.Name == p.currentLabelPair.Name {
			p.parseError(fmt.Sprintf("duplicate label names for metric %q", p.metric.Name))
			return nil
		}
	}
	return p.startLabelValue
}

// startLabelValue represents the state where the next byte read from p.buf is
// the start of a (quoted) label value (or whitespace leading up to it).
func (p *TextParser) startLabelValue() stateFn {
	if p.skipBlankTab(); p.err != nil {
		return nil // Unexpected end of input.
	}
	if p.currentByte != '"' {
		p.parseError(fmt.Sprintf("expected '\"' at start of label value, found %q", p.currentByte))
		return nil
	}
	if p.readTokenAsLabelValue(); p.err != nil {
		return nil
	}
	if !utf8.ValidString(p.currentToken.String()) {
		p.parseError(fmt.Sprintf("invalid label value %q", p.currentToken.String()))
		return nil
	}
	p.currentLabelPair.Value = p.currentToken.String()
	p.metric.Labels = append(p.metric.Labels, *p.currentLabelPair)
	if p.skipBlankTab(); p.err != nil {
		return nil // Unexpected end of input.
	}
	switch p.currentByte {
	case ',':
		return p.startLabelName

	case '}':
		if p.skipBlankTab(); p.err != nil {
			return nil // Unexpected end of input.
		}
		return p.readingValue
	default:
		p.parseError(fmt.Sprintf("unexpected end of label value %q", p.currentLabelPair.Value))
		return nil
	}
}

// readingValue represents the state where the last byte read (now in
// p.currentByte) is the first byte of the sample value (i.e. a float).
func (p *TextParser) readingValue() stateFn {
	if p.readTokenUntilWhitespace(); p.err != nil {
		// Metric line might end with value
		if p.err == io.EOF {
			p.err = nil
		} else {
			return nil // Unexpected end of input.
		}
	}
	value, err := parseFloat(p.currentToken.String())
	if err != nil {
		// Create a more helpful error message.
		p.parseError(fmt.Sprintf("expected float as value, got %q", p.currentToken.String()))
		return nil
	}
	p.metric.Value = value
	if p.currentByte == '\n' {
		return nil
	}
	return p.startTimestamp
}

// startTimestamp represents the state where the next byte read from p.buf is
// the start of the timestamp (or whitespace leading up to it).
func (p *TextParser) startTimestamp() stateFn {
	if p.skipBlankTab(); p.err != nil {
		return nil // Unexpected end of input.
	}
	if p.readTokenUntilWhitespace(); p.err != nil {
		// Metric line might end with timestamp
		if p.err == io.EOF {
			p.err = nil
		} else {
			return nil // Unexpected end of input.
		}
	}
	timestamp, err := strconv.ParseInt(p.currentToken.String(), 10, 64)
	if err != nil {
		// Create a more helpful error message.
		p.parseError(fmt.Sprintf("expected integer as timestamp, got %q", p.currentToken.String()))
		return nil
	}
	p.metric.TimestampMs = timestamp
	if p.readTokenUntilNewline(false); p.err != nil {
		return nil // Unexpected end of input.
	}
	if p.currentToken.Len() > 0 {
		p.parseError(fmt.Sprintf("spurious string after timestamp: %q", p.currentToken.String()))
		return nil
	}
	return p.startOfLine
}

// parseError sets p.err to a ParseError at the current line with the given
// message.
func (p *TextParser) parseError(msg string) {
	p.err = ParseError{
		Msg: msg,
	}
}

// skipBlankTab reads (and discards) bytes from p.buf until it encounters a byte
// that is neither ' ' nor '\t'. That byte is left in p.currentByte.
func (p *TextParser) skipBlankTab() {
	for {
		if p.currentByte, p.err = p.buf.ReadByte(); p.err != nil || !isBlankOrTab(p.currentByte) {
			return
		}
	}
}

// skipBlankTabIfCurrentBlankTab works exactly as skipBlankTab but doesn't do
// anything if p.currentByte is neither ' ' nor '\t'.
func (p *TextParser) skipBlankTabIfCurrentBlankTab() {
	if isBlankOrTab(p.currentByte) {
		p.skipBlankTab()
	}
}

// readTokenUntilWhitespace copies bytes from p.buf into p.currentToken.  The
// first byte considered is the byte already read (now in p.currentByte).  The
// first whitespace byte encountered is still copied into p.currentByte, but not
// into p.currentToken.
func (p *TextParser) readTokenUntilWhitespace() {
	p.currentToken.Reset()
	for p.err == nil && !isBlankOrTab(p.currentByte) && p.currentByte != '\n' {
		p.currentToken.WriteByte(p.currentByte)
		p.currentByte, p.err = p.buf.ReadByte()
	}
}

// readTokenUntilNewline copies bytes from p.buf into p.currentToken.  The first
// byte considered is the byte already read (now in p.currentByte).  The first
// newline byte encountered is still copied into p.currentByte, but not into
// p.currentToken. If recognizeEscapeSequence is true, two escape sequences are
// recognized: '\\' translates into '\', and '\n' into a line-feed character.
// All other escape sequences are invalid and cause an error.
func (p *TextParser) readTokenUntilNewline(recognizeEscapeSequence bool) {
	p.currentToken.Reset()
	escaped := false
	for p.err == nil {
		if recognizeEscapeSequence && escaped {
			switch p.currentByte {
			case '\\':
				p.currentToken.WriteByte(p.currentByte)
			case 'n':
				p.currentToken.WriteByte('\n')
			default:
				p.parseError(fmt.Sprintf("invalid escape sequence '\\%c'", p.currentByte))
				return
			}
			escaped = false
		} else {
			switch p.currentByte {
			case '\n':
				return
			case '\\':
				escaped = true
			default:
				p.currentToken.WriteByte(p.currentByte)
			}
		}
		p.currentByte, p.err = p.buf.ReadByte()
	}
}

// readTokenAsMetricName copies a metric name from p.buf into p.currentToken.
// The first byte considered is the byte already read (now in p.currentByte).
// The first byte not part of a metric name is still copied into p.currentByte,
// but not into p.currentToken.
func (p *TextParser) readTokenAsMetricName() {
	p.currentToken.Reset()
	if !isValidMetricNameStart(p.currentByte) {
		fmt.Printf("INVALID BYTE %v", p.currentByte)
		return
	}
	for {
		p.currentToken.WriteByte(p.currentByte)
		p.currentByte, p.err = p.buf.ReadByte()
		if p.err != nil || !isValidMetricNameContinuation(p.currentByte) {
			return
		}
	}
}

// readTokenAsLabelName copies a label name from p.buf into p.currentToken.
// The first byte considered is the byte already read (now in p.currentByte).
// The first byte not part of a label name is still copied into p.currentByte,
// but not into p.currentToken.
func (p *TextParser) readTokenAsLabelName() {
	p.currentToken.Reset()
	if !isValidLabelNameStart(p.currentByte) {
		return
	}
	for {
		p.currentToken.WriteByte(p.currentByte)
		p.currentByte, p.err = p.buf.ReadByte()
		if p.err != nil || !isValidLabelNameContinuation(p.currentByte) {
			return
		}
	}
}

// readTokenAsLabelValue copies a label value from p.buf into p.currentToken.
// In contrast to the other 'readTokenAs...' functions, which start with the
// last read byte in p.currentByte, this method ignores p.currentByte and starts
// with reading a new byte from p.buf. The first byte not part of a label value
// is still copied into p.currentByte, but not into p.currentToken.
func (p *TextParser) readTokenAsLabelValue() {
	p.currentToken.Reset()
	escaped := false
	for {
		if p.currentByte, p.err = p.buf.ReadByte(); p.err != nil {
			return
		}
		if escaped {
			switch p.currentByte {
			case '"', '\\':
				p.currentToken.WriteByte(p.currentByte)
			case 'n':
				p.currentToken.WriteByte('\n')
			default:
				p.parseError(fmt.Sprintf("invalid escape sequence '\\%c'", p.currentByte))
				return
			}
			escaped = false
			continue
		}
		switch p.currentByte {
		case '"':
			return
		case '\n':
			p.parseError(fmt.Sprintf("label value %q contains unescaped new-line", p.currentToken.String()))
			return
		case '\\':
			escaped = true
		default:
			p.currentToken.WriteByte(p.currentByte)
		}
	}
}

func isValidLabelNameStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isValidLabelNameContinuation(b byte) bool {
	return isValidLabelNameStart(b) || (b >= '0' && b <= '9')
}

func isValidMetricNameStart(b byte) bool {
	return isValidLabelNameStart(b) || b == ':'
}

func isValidMetricNameContinuation(b byte) bool {
	return isValidLabelNameContinuation(b) || b == ':'
}

func isBlankOrTab(b byte) bool {
	return b == ' ' || b == '\t'
}

func parseFloat(s string) (float64, error) {
	if strings.ContainsAny(s, "pP_") {
		return 0, fmt.Errorf("unsupported character in float")
	}
	return strconv.ParseFloat(s, 64)
}
