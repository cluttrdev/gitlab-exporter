package rest

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"unicode/utf8"

	_gitlab "gitlab.com/gitlab-org/api/client-go"

	"go.cluttr.dev/gitlab-exporter/internal/expfmt"
)

type JobLogData struct {
	Sections   []SectionData  `json:"sections"`
	Metrics    []MetricData   `json:"metrics"`
	Properties []PropertyData `json:"properties"`
}

type MetricData struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
}

type PropertyData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (c *Client) GetJobLogData(ctx context.Context, projectId int64, jobId int64) (JobLogData, error) {
	log, _, err := c.client.Jobs.GetTraceFile(int(projectId), int(jobId), _gitlab.WithContext(ctx))
	if err != nil {
		return JobLogData{}, fmt.Errorf("get job log: %w", err)
	}

	data, err := ParseJobLog(log)
	if err != nil {
		return JobLogData{}, fmt.Errorf("parse job log: %w", err)
	}

	return data, nil
}

func ParseJobLog(trace *bytes.Reader) (JobLogData, error) {
	var (
		data           JobLogData
		sections       sectionStack
		parser         expfmt.TextParser
		propertyParser jobLogPropertyParser
	)

	const (
		METRIC_MARKER   = `METRIC_`
		PROPERTY_MARKER = `PROPERTY_`
		SECTION_MARKER  = `section_`
	)

	scanner := bufio.NewScanner(trace)
	for scanner.Scan() {
		line := scanner.Bytes()

		if bytes.HasPrefix(line, []byte(METRIC_MARKER)) {
			metric, err := parser.LineToMetric(line[len(METRIC_MARKER):])
			if err != nil || metric == nil {
				// we ignore parsing errors here
				// TODO: should we handle them somehow?
				continue
			}

			data.Metrics = append(data.Metrics, MetricData{
				Name:      metric.Name,
				Labels:    convertMetricLabels(metric.Labels),
				Value:     metric.Value,
				Timestamp: metric.TimestampMs,
			})
		}
		if bytes.HasPrefix(line, []byte(PROPERTY_MARKER)) {
			property, err := propertyParser.LineToProperty(line[len(PROPERTY_MARKER):])
			if err != nil || property == nil {
				// we ignore parsing errors here
				// TODO: should we handle them somehow?
				continue
			}

			data.Properties = append(data.Properties, *property)
		}

		var i, j int
		sep := []byte(SECTION_MARKER)
		for {
			j = bytes.Index(line[i:], sep)
			if j < 0 {
				break
			}

			marker, ts, name, err := parseSection(line[i:])
			if err != nil {
				// TODO: what?
			} else if marker == string(sectionMarkerStart) {
				sections.Start(ts, name)
			} else if marker == string(sectionMarkerEnd) {
				data.Sections = append(data.Sections, sections.End(ts, name)...)
			}

			i = i + j + 1
		}
	}

	// add all unfinished sections (e.g. due to job interruption) open-ended
	// to give caller a chance to set end timestamp
	for sections.Size() > 0 {
		data.Sections = append(data.Sections, sections.Pop())
	}

	return data, nil
}

func convertMetricLabels(pairs []expfmt.MetricLabelPair) map[string]string {
	if len(pairs) == 0 {
		return nil
	}
	labels := map[string]string{}
	for _, p := range pairs {
		labels[p.Name] = p.Value
	}
	return labels
}

type jobLogPropertyParser struct {
	buf *bufio.Reader

	currentByte  byte
	currentToken bytes.Buffer
	err          error

	labelPair PropertyData
}

func (p *jobLogPropertyParser) LineToProperty(line []byte) (*PropertyData, error) {
	if len(line) == 0 {
		return nil, nil
	}
	p.reset(bytes.NewReader(line))

	return p.lineToProperty()
}

func (p *jobLogPropertyParser) reset(in io.Reader) {
	p.labelPair = PropertyData{}
	if p.buf == nil {
		p.buf = bufio.NewReader(in)
	} else {
		p.buf.Reset(in)
	}
	p.err = nil
}

func (p *jobLogPropertyParser) lineToProperty() (*PropertyData, error) {
	p.currentByte, p.err = p.buf.ReadByte()
	if p.err != nil {
		return nil, p.err
	}

	// parse name
	if p.readTokenAsLabelName(); p.err != nil {
		return nil, p.err
	}
	if p.currentToken.Len() == 0 {
		return nil, fmt.Errorf("empty label name")
	}
	p.labelPair.Name = p.currentToken.String()

	// check for '='
	if p.skipBlankTabIfCurrentBlankTab(); p.err != nil {
		return nil, p.err
	}
	if p.currentByte != '=' {
		return nil, fmt.Errorf("expected '=' after label name, found %q", p.currentByte)
	}

	// parse value
	if p.skipBlankTab(); p.err != nil {
		return nil, p.err
	}
	if p.currentByte != '"' {
		return nil, fmt.Errorf("expected '\"' at start of label value, found %q", p.currentByte)
	}
	if p.readTokenAsLabelValue(); p.err != nil {
		return nil, p.err
	}

	value := p.currentToken.String()
	if !utf8.ValidString(value) {
		return nil, fmt.Errorf("invalid label value %q", value)
	}
	p.labelPair.Value = value

	return &p.labelPair, nil
}

func (p *jobLogPropertyParser) readTokenAsLabelName() {
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

func (p *jobLogPropertyParser) readTokenAsLabelValue() {
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
				p.err = fmt.Errorf("invalid escape sequence '\\%c'", p.currentByte)
				return
			}
			escaped = false
			continue
		}
		switch p.currentByte {
		case '"':
			return
		case '\n':
			p.err = fmt.Errorf("label value %q contains unescaped new-line", p.currentToken.String())
			return
		case '\\':
			escaped = true
		default:
			p.currentToken.WriteByte(p.currentByte)
		}
	}
}

func (p *jobLogPropertyParser) skipBlankTabIfCurrentBlankTab() {
	if isBlankOrTab(p.currentByte) {
		p.skipBlankTab()
	}
}

func (p *jobLogPropertyParser) skipBlankTab() {
	for {
		if p.currentByte, p.err = p.buf.ReadByte(); p.err != nil || !isBlankOrTab(p.currentByte) {
			return
		}
	}
}

func isValidLabelNameStart(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || b == '_'
}

func isValidLabelNameContinuation(b byte) bool {
	return isValidLabelNameStart(b) || (b >= '0' && b <= '9')
}

func isBlankOrTab(b byte) bool {
	return b == ' ' || b == '\t'
}
