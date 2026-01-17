package parser

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

type HTTPRequest struct {
	startLine string
	headers   []string
	body      string
}

func Listen(conn net.Conn) *Parser {
	return &Parser{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
}

func (p *Parser) Parse() (*HTTPRequest, error) {
	var request HTTPRequest
	startLine, err := p.parseStartline()
	if err != nil {
		return nil, fmt.Errorf("failed parsing startline: %s", err)
	}
	request.startLine = startLine

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
