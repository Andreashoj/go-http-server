package router

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Parser struct {
	conn   net.Conn
	reader *bufio.Reader
}

func Listen(conn net.Conn) *Parser {
	return &Parser{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
}

func (p *Parser) Parse() (HTTPRequest, error) {
	var request httpRequest
	startLine, err := p.parseStartline()

	if err != nil {
		return nil, fmt.Errorf("failed parsing startline: %s", err)
	}
	request.startLine = startLine
	request.url = p.parseUrl(startLine)
	request.params = p.parseParams(startLine)
	request.method = p.parseMethod(startLine)

	headers, err := p.parseHeaders()
	if err != nil {
		return nil, fmt.Errorf("failed parsing headers: %s", err)
	}

	request.headers = headers
	contentLength := p.getContentLength(headers)

	// No body, return early
	if contentLength == 0 {
		return &request, nil
	}

	body, err := p.parseBody(contentLength)
	if err != nil {
		return nil, fmt.Errorf("failed parsing body: %s", err)
	}

	request.body = body
	return &request, nil
}

func (p *Parser) parseStartline() (string, error) {
	startLine, err := p.reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed decoding starting, err: %s", err)
	}

	return startLine, nil
}

func (p *Parser) parseMethod(startLine string) Request {
	return Request(strings.Split(startLine, " ")[0])
}

func (p *Parser) parseUrl(startLine string) string {
	url := strings.Split(startLine, " ")[1]

	i := strings.Index(url, "?")
	if i == -1 { // NO query params present
		return url
	}

	separatedUrl := strings.Split(url, "?")
	return separatedUrl[0]
}

func (p *Parser) parseParams(startLine string) map[string]string {
	params := make(map[string]string)
	url := strings.Split(startLine, " ")[1]

	i := strings.Index(url, "?")
	if i == -1 { // NO query params present
		return params
	}

	// ? value=1&another-val=23 ?
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

func (p *Parser) parseHeaders() ([]string, error) {
	var headers []string
	for {
		line, err := p.reader.ReadString('\n')
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

func (p *Parser) getContentLength(headers []string) int {
	for _, header := range headers {
		h := strings.Split(header, ":")
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

func (p *Parser) parseBody(contentLength int) (string, error) {
	body := make([]byte, contentLength)
	_, err := p.reader.Read(body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
