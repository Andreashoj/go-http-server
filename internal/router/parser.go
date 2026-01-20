package router

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Parse(reader *bufio.Reader) (HTTPRequest, error) {
	var request httpRequest

	startLine, err := parseStartline(reader)
	if err != nil {
		fmt.Println("here", err)
		return &request, fmt.Errorf("failed parsing startline: %s", err)
	}

	request.startLine = startLine
	request.url = parseUrl(startLine)
	request.params = parseParams(startLine)
	request.method = parseMethod(startLine)

	headers, err := parseHeaders(reader)
	if err != nil {
		return nil, fmt.Errorf("failed parsing headers: %s", err)
	}

	request.headers = headers
	contentLength := getContentLength(headers)

	body, err := parseBody(reader, contentLength)
	if err != nil {
		return nil, fmt.Errorf("failed parsing body: %s", err)
	}

	fmt.Println("here")

	request.body = body
	return &request, nil
}

func parseStartline(reader *bufio.Reader) (string, error) {
	startLine, err := reader.ReadString('\n')
	if err != nil {
		return "", err
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

func parseParams(startLine string) map[string]string {
	params := make(map[string]string)
	url := strings.Split(startLine, " ")[1]

	i := strings.Index(url, "?")
	if i == -1 { // NO query params present
		return params
	}

	query := strings.Split(url, "?")
	if query[len(query)-1] == "?" { // If last letter is ?, there is no params in the URL => google.com?
		return params
	}

	pr := strings.Split(query[len(query)-1], "&")
	for _, entry := range pr {
		parts := strings.Split(entry, "=")
		if len(parts) != 2 {
			fmt.Printf("failed decoding parameter on input: %s", entry)
			continue
		}
		key, value := parts[0], parts[1]
		params[key] = value
	}

	return params
}

func parseHeaders(reader *bufio.Reader) ([]string, error) {
	var headers []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed decoding header, err: %s", err)
		}

		// body starts
		if line == "\r\n" {
			break
		}

		headers = append(headers, line)
	}

	return headers, nil
}

func getContentLength(headers []string) int {
	for _, hder := range headers {
		h := strings.Split(hder, ":")
		if h[0] == "Content-Length" {
			length, err := strconv.Atoi(strings.TrimSpace(h[1]))
			if err != nil {
				fmt.Printf("failed retrieving content length from header: %s", err)
				return 0
			}

			return length
		}
	}

	return 0
}

func parseBody(reader *bufio.Reader, contentLength int) (string, error) {
	body := make([]byte, contentLength)
	n, err := io.ReadFull(reader, body)
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return "", err
	}

	return string(body[:n]), nil
}
