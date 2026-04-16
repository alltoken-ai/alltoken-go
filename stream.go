package alltoken

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"
)

// Stream reads SSE events from a streaming chat completion response.
type Stream struct {
	resp    *http.Response
	scanner *bufio.Scanner
	done    bool
	err     error
}

func newStream(resp *http.Response) *Stream {
	return &Stream{
		resp:    resp,
		scanner: bufio.NewScanner(resp.Body),
	}
}

// Next advances to the next chunk. Returns false when the stream is done or errored.
func (s *Stream) Next() bool {
	if s.done {
		return false
	}
	for s.scanner.Scan() {
		line := s.scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if strings.TrimSpace(data) == "[DONE]" {
				s.done = true
				return false
			}
			return true
		}
	}
	s.err = s.scanner.Err()
	s.done = true
	return false
}

// Current returns the current chunk. Call after Next() returns true.
func (s *Stream) Current() (*ChatCompletionChunk, error) {
	line := s.scanner.Text()
	data := strings.TrimPrefix(line, "data: ")
	var chunk ChatCompletionChunk
	if err := json.Unmarshal([]byte(data), &chunk); err != nil {
		return nil, err
	}
	return &chunk, nil
}

// Close closes the underlying response body.
func (s *Stream) Close() error {
	return s.resp.Body.Close()
}

// Err returns any error encountered during iteration.
func (s *Stream) Err() error {
	return s.err
}
