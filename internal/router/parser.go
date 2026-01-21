package router

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

func Parse(reader *bufio.Reader) (HTTPRequest, error) {
	var request httpRequest

	// Handle startline
	startLine, err := parseStartline(reader)
	if err != nil {
		return &request, fmt.Errorf("failed parsing startline: %s", err)
	}

	request.startLine = startLine
	request.method = parseMethod(startLine)
	request.url = parseUrl(startLine)
	params, err := parseParams(startLine)
	if err != nil {
		return nil, fmt.Errorf("failed parsing params: %s", err)
	}
	request.params = params

	// Handle headers
	headers, err := parseHeaders(reader)
	if err != nil {
		return nil, fmt.Errorf("failed parsing headers: %s", err)
	}
	request.headers = headers
	contentLength, err := getContentLength(headers)
	if err != nil {
		return nil, fmt.Errorf("content length is specified but failed retrieving it: %s", err)
	}

	// Handle body
	body, err := parseBody(reader, contentLength)
	if err != nil {
		return nil, fmt.Errorf("failed parsing body: %s", err)
	}

	request.body = body
	return &request, nil
}

func parseStartline(reader *bufio.Reader) (string, error) {
	startLine, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	if len(strings.Split(startLine, " ")) != 3 {
		return "", fmt.Errorf("expected startline to have method, url and http version. One or more is missing: %s", startLine)
	}

	return startLine, nil
}

func parseMethod(startLine string) Request {
	return Request(strings.Split(startLine, " ")[0])
}

func parseUrl(startLine string) string {
	url := strings.Split(startLine, " ")[1]

	i := strings.Index(url, "?")
	if i == -1 { // NO query params present
		return url
	}

	separatedUrl := strings.Split(url, "?")
	return separatedUrl[0]
}

func parseParams(startLine string) (map[string]string, error) {
	params := make(map[string]string)
	endpoint := strings.Split(startLine, " ")[1]

	i := strings.Index(endpoint, "?")
	if i == -1 { // NO query params present
		return params, nil
	}

	query := strings.Split(endpoint, "?")
	if query[len(query)-1] == "?" { // If last letter is ?, there is no params in the URL => google.com?
		return params, nil
	}

	pr := strings.Split(query[len(query)-1], "&")
	for _, entry := range pr {
		parts := strings.Split(entry, "=")
		if len(parts) != 2 {
			fmt.Printf("failed decoding parameter on input: %s", entry)
			continue
		}
		key, value := parts[0], parts[1]
		decoded, err := url.QueryUnescape(value)
		if err != nil {
			return nil, fmt.Errorf("failed decoding parameter value for key: %s with value: %s and error: %s", key, value, err)
		}

		params[key] = decoded
	}

	return params, nil
}

func parseHeaders(reader *bufio.Reader) ([]string, error) {
	var headers []string
	var hasHost bool
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed decoding header, err: %s", err)
		}

		// body starts
		if line == "\r\n" {
			break
		}

		// Validate "value: key" format
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed header format")
		}

		headerKey := strings.TrimSpace(strings.ToLower(string(parts[0])))
		headerValue := strings.TrimSpace(strings.ToLower(string(parts[1])))

		// Validate host exists
		if headerKey == "host" && headerValue != "" {
			hasHost = true
		}

		// Validate content length value
		if headerKey == "content-length" {
			length, err := strconv.Atoi(headerValue)
			if err != nil || length < 0 {
				return nil, fmt.Errorf("invalid content length value: %s", err)
			}
		}

		headers = append(headers, line)
	}

	if !hasHost {
		return nil, fmt.Errorf("failed because no host header was present")
	}

	return headers, nil
}

func getContentLength(headers []string) (int, error) {
	for _, hder := range headers {
		h := strings.Split(hder, ":")
		if h[0] == "Content-Length" {
			length, _ := strconv.Atoi(strings.TrimSpace(h[1])) // validation of content length handled in header parser
			return length, nil
		}
	}

	return 0, nil
}

func parseBody(reader *bufio.Reader, contentLength int) (string, error) {
	body := make([]byte, contentLength)
	n, err := io.ReadFull(reader, body)
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}

	return string(body[:n]), nil
}
