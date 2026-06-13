package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
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

type EmptyLine struct{}

type TransactionHeaderLine struct {
	Date        string
	Description string
}

type TransactionDetailLine struct {
	Account string
	Amount  float64
}

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

func parseLine(line string, lineNumber int) (any, error) {
	if strings.TrimSpace(line) == "" {
		return &EmptyLine{}, nil
	}

	if result, err := parseTransactionHeaderLine(line); err == nil {
		return result, nil
	}

	if result, err := parseTransactionDetailLine(line); err == nil {
		return result, nil
	}

	return nil, fmt.Errorf("invalid line: %s at line number %d", line, lineNumber)
}

// parses the line `2026-01-01  Initial balance` and returns a TransactionHeaderLine
func parseTransactionHeaderLine(line string) (*TransactionHeaderLine, error) {
	regex := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2})\s+(.*)$`)
	matches := regex.FindStringSubmatch(line)

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid transaction header line: %s", line)
	}

	dateString := matches[1]
	description := matches[2]

	return &TransactionHeaderLine{
		Date:        dateString,
		Description: description,
	}, nil
}

// parses the line `  Assets:Cash in Bank  1000.00` and returns a TransactionDetailLine
func parseTransactionDetailLine(line string) (*TransactionDetailLine, error) {
	regex := regexp.MustCompile(`^\s+(\w+[:\w]*(?:\s\w+[:\w]*)*)\s{2,}([-+]?\d*\.?\d+)$`)
	matches := regex.FindStringSubmatch(line)

	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid transaction detail line: %s", line)
	}

	account := matches[1]
	amount := matches[2]

	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid transaction detail line: %s", line)
	}

	return &TransactionDetailLine{
		Account: account,
		Amount:  amountFloat,
	}, nil
}
