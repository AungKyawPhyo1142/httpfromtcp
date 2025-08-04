package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("could not read from reader: %w", err)
	}

	requestLine, err := parseRequestLine(b)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}

// parse the request-line from TCP, to HTTP request-line
func parseRequestLine(b []byte) (RequestLine, error) {
	rawRequest := string(b)

	// in HTTP, newlines are \r\n. We only care about the first line.
	firstLine, _, _ := strings.Cut(rawRequest, "\r\n")

	// The request line has 3 parts: METHOD /path HTTP/VERSION
	parts := strings.Split(firstLine, " ")
	if len(parts) != 3 {
		return RequestLine{}, fmt.Errorf("malformed request line: expected 3 parts, got %d", len(parts))
	}

	method, target, version := parts[0], parts[1], parts[2]

	// Verify that the "method" part only contains capital alphabetic characters.
	for _, char := range method {
		if char < 'A' || char > 'Z' {
			return RequestLine{}, fmt.Errorf("invalid method '%s': must be all uppercase", method)
		}
	}

	// Verify that the http version part is 1.1
	if version != "HTTP/1.1" {
		return RequestLine{}, fmt.Errorf("invalid http version '%s': only HTTP/1.1 is supported", version)
	}

	return RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, nil
}
