package parser

import (
	"bufio"
	"io"
	"strings"
)

/*
This package is used to parse the input file and return the data in a structured format.

Ledger file format:

	2026-01-01  Initial balance
		Assets:Cash    1000.00
		Equity:Opening Balances    -1000.00

	2026-01-02 Rent
		Expenses:Rent    1000.00
		Assets:Cash    -1000.00
*/

type Parser struct {
	reader io.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		reader: reader,
	}
}

func (p *Parser) Parse() ([]any, error) {
	scanner := bufio.NewScanner(p.reader)
	scanner.Split(bufio.ScanLines)
	lineNumber := 0
	tokens := make([]any, 0)

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimRight(scanner.Text(), " \n")
		result, err := parseLine(line, lineNumber)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, result)
	}

	return tokens, nil
}
