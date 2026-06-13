package parser

import (
	"strings"
	"testing"
)

func TestParser_Parse(t *testing.T) {
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
		tokens, err := parser.Parse()
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
		tokens, err := parser.Parse()
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
