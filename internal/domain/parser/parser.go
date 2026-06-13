package parser

import (
	"bufio"
	"fmt"
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

type Transaction struct {
	Date        string
	Description string
	Lines       []TransactionLine
}

type TransactionLine struct {
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

func (p *Parser) ParseTokens() ([]any, error) {
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

func (p *Parser) ParseTransactions() ([]*Transaction, error) {

	tokens, err := p.ParseTokens()

	if err != nil {
		return nil, err
	}

	stack := make([]any, 0)
	transactions := make([]*Transaction, 0)

	var currentTransaction *Transaction

	for _, token := range tokens {
		switch token := token.(type) {
		case *TransactionHeaderLine:
			currentTransaction = &Transaction{
				Date:        token.Date,
				Description: token.Description,
				Lines:       make([]TransactionLine, 0),
			}
			stack = append(stack, token)
		case *TransactionDetailLine:
			peek := stack[len(stack)-1]
			if _, ok := peek.(*TransactionHeaderLine); !ok {
				return nil, fmt.Errorf("invalid transaction detail line: %v", token)
			}
			currentTransaction.Lines = append(currentTransaction.Lines, TransactionLine{
				Account: token.Account,
				Amount:  token.Amount,
			})
		case *EmptyLine:
			peek := stack[len(stack)-1]
			if _, ok := peek.(*TransactionHeaderLine); ok {
				transactions = append(transactions, currentTransaction)
				currentTransaction = nil
				stack = stack[:len(stack)-1]
			}
		}
	}

	return transactions, nil
}
