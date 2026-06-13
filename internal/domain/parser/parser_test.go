package parser

import (
	"strings"
	"testing"
)

func TestParser_ParseTokens(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		input := strings.Join([]string{
			"2026-01-01  Initial balance",
			"\tAssets:Cash    1000.00",
			"\tEquity:Opening Balances    -1000.00",
			"",
			"2026-01-02 Rent",
			"\tExpenses:Rent    1000.00",
			"\tAssets:Cash    -1000.00",
		}, "\n")

		parser := NewParser(strings.NewReader(input))
		tokens, err := parser.ParseTokens()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		want := []any{
			&TransactionHeaderLine{Date: "2026-01-01", Description: "Initial balance"},
			&TransactionDetailLine{Account: "Assets:Cash", Amount: 1000.00},
			&TransactionDetailLine{Account: "Equity:Opening Balances", Amount: -1000.00},
			&EmptyLine{},
			&TransactionHeaderLine{Date: "2026-01-02", Description: "Rent"},
			&TransactionDetailLine{Account: "Expenses:Rent", Amount: 1000.00},
			&TransactionDetailLine{Account: "Assets:Cash", Amount: -1000.00},
		}

		if len(tokens) != len(want) {
			t.Fatalf("expected %d tokens, got %d", len(want), len(tokens))
		}
		for i := range want {
			assertToken(t, tokens[i], want[i])
		}
	})

	t.Run("error path", func(t *testing.T) {
		input := strings.Join([]string{
			"2026-01-01  Initial balance",
			"\tAssets:Cash    1000.00",
			"invalid line here",
		}, "\n")

		parser := NewParser(strings.NewReader(input))
		tokens, err := parser.ParseTokens()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if tokens != nil {
			t.Fatalf("expected nil tokens, got %v", tokens)
		}
		if !strings.Contains(err.Error(), "invalid line: invalid line here at line number 3") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestParser_ParseTransactions(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		input := strings.Join([]string{
			"2026-01-01  Initial balance",
			"\tAssets:Cash    1000.00",
			"\tEquity:Opening Balances    -1000.00",
			"",
			"2026-01-02 Rent",
			"\tExpenses:Rent    1000.00",
			"\tAssets:Cash    -1000.00",
			"",
			"",
		}, "\n")

		parser := NewParser(strings.NewReader(input))
		transactions, err := parser.ParseTransactions()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		want := []*Transaction{
			{
				Date:        "2026-01-01",
				Description: "Initial balance",
				Lines: []TransactionLine{
					{Account: "Assets:Cash", Amount: 1000.00},
					{Account: "Equity:Opening Balances", Amount: -1000.00},
				},
			},
			{
				Date:        "2026-01-02",
				Description: "Rent",
				Lines: []TransactionLine{
					{Account: "Expenses:Rent", Amount: 1000.00},
					{Account: "Assets:Cash", Amount: -1000.00},
				},
			},
		}

		if len(transactions) != len(want) {
			t.Fatalf("expected %d transactions, got %d", len(want), len(transactions))
		}
		for i := range want {
			assertTransaction(t, transactions[i], want[i])
		}
	})

	t.Run("error path", func(t *testing.T) {
		input := strings.Join([]string{
			"2026-01-01  Initial balance",
			"\tAssets:Cash    1000.00",
			"invalid line here",
		}, "\n")

		parser := NewParser(strings.NewReader(input))
		transactions, err := parser.ParseTransactions()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if transactions != nil {
			t.Fatalf("expected nil transactions, got %v", transactions)
		}
		if !strings.Contains(err.Error(), "invalid line: invalid line here at line number 3") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func assertTransaction(t *testing.T, got, want *Transaction) {
	t.Helper()

	if got.Date != want.Date {
		t.Errorf("expected date %s, got %s", want.Date, got.Date)
	}
	if got.Description != want.Description {
		t.Errorf("expected description %s, got %s", want.Description, got.Description)
	}
	if len(got.Lines) != len(want.Lines) {
		t.Fatalf("expected %d lines, got %d", len(want.Lines), len(got.Lines))
	}
	for i := range want.Lines {
		if got.Lines[i].Account != want.Lines[i].Account {
			t.Errorf("line %d: expected account %s, got %s", i, want.Lines[i].Account, got.Lines[i].Account)
		}
		if got.Lines[i].Amount != want.Lines[i].Amount {
			t.Errorf("line %d: expected amount %f, got %f", i, want.Lines[i].Amount, got.Lines[i].Amount)
		}
	}
}
