package parser

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseTransactionHeaderLine(t *testing.T) {
	tests := []struct {
		line          string
		expected      *TransactionHeaderLine
		expectedError error
	}{
		{
			line: "2026-01-01  Initial balance",
			expected: &TransactionHeaderLine{
				Date:        "2026-01-01",
				Description: "Initial balance",
			},
			expectedError: nil,
		},
		{
			line: "2026-01-01 Rent",
			expected: &TransactionHeaderLine{
				Date:        "2026-01-01",
				Description: "Rent",
			},
			expectedError: nil,
		},
		{
			line: "2026-01-01Rent",
			expected: &TransactionHeaderLine{
				Date:        "2026-01-01",
				Description: "Rent",
			},
			expectedError: fmt.Errorf("invalid transaction header line: %s", "2026-01-01Rent"),
		},
	}

	for _, test := range tests {
		result, err := parseTransactionHeaderLine(test.line)
		if err != nil && test.expectedError == nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err == nil && test.expectedError != nil {
			t.Errorf("expected error %v, got nil", test.expectedError)
		}
		if err != nil && test.expectedError != nil {
			if err.Error() != test.expectedError.Error() {
				t.Errorf("expected error %v, got %v", test.expectedError, err)
			}
		}
		if result != nil && test.expected != nil {
			if result.Date != test.expected.Date {
				t.Errorf("expected date %s, got %s", test.expected.Date, result.Date)
			}
			if result.Description != test.expected.Description {
				t.Errorf("expected description %s, got %s", test.expected.Description, result.Description)
			}
		}
	}

}

func TestParseTransactionDetailLine(t *testing.T) {
	tests := []struct {
		line          string
		expected      *TransactionDetailLine
		expectedError error
	}{
		{
			line: "  Assets:Cash  1000.00",
			expected: &TransactionDetailLine{
				Account: "Assets:Cash",
				Amount:  1000.00,
			},
			expectedError: nil,
		},
		{
			line: "  Assets:Cash in Bank  1000.00",
			expected: &TransactionDetailLine{
				Account: "Assets:Cash in Bank",
				Amount:  1000.00,
			},
			expectedError: nil,
		},
		{
			line:          "  Assets:Cash in  Bank  1000.00",
			expected:      nil,
			expectedError: fmt.Errorf("invalid transaction detail line: %s", "  Assets:Cash in  Bank  1000.00"),
		},
		{
			line: "  Equity:Opening Balances    -1000.00",
			expected: &TransactionDetailLine{
				Account: "Equity:Opening Balances",
				Amount:  -1000.00,
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		result, err := parseTransactionDetailLine(test.line)
		if err != nil && test.expectedError == nil {
			t.Errorf("expected no error, got %v", err)
		}
		if err == nil && test.expectedError != nil {
			t.Errorf("expected error %v, got nil", test.expectedError)
		}
		if err != nil && test.expectedError != nil {
			if err.Error() != test.expectedError.Error() {
				t.Errorf("expected error %v, got %v", test.expectedError, err)
			}
		}
		if result != nil && test.expected != nil {
			if result.Account != test.expected.Account {
				t.Errorf("expected account %s, got %s", test.expected.Account, result.Account)
			}
			if result.Amount != test.expected.Amount {
				t.Errorf("expected amount %f, got %f", test.expected.Amount, result.Amount)
			}
		}
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		line          string
		lineNumber    int
		expected      any
		expectedError string
	}{
		{
			line:       "",
			lineNumber: 1,
			expected:   &EmptyLine{},
		},
		{
			line:       "   ",
			lineNumber: 2,
			expected:   &EmptyLine{},
		},
		{
			line:       "2026-01-01  Initial balance",
			lineNumber: 3,
			expected: &TransactionHeaderLine{
				Date:        "2026-01-01",
				Description: "Initial balance",
			},
		},
		{
			line:       "  Assets:Cash  1000.00",
			lineNumber: 4,
			expected: &TransactionDetailLine{
				Account: "Assets:Cash",
				Amount:  1000.00,
			},
		},
		{
			line:          "not a ledger line",
			lineNumber:    5,
			expectedError: "invalid line: not a ledger line at line number 5",
		},
	}

	for _, test := range tests {
		result, err := parseLine(test.line, test.lineNumber)
		if test.expectedError != "" {
			if err == nil {
				t.Errorf("line %q: expected error, got nil", test.line)
				continue
			}
			if err.Error() != test.expectedError {
				t.Errorf("line %q: expected error %q, got %q", test.line, test.expectedError, err.Error())
			}
			continue
		}
		if err != nil {
			t.Errorf("line %q: expected no error, got %v", test.line, err)
			continue
		}
		assertToken(t, result, test.expected)
	}
}

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

func assertToken(t *testing.T, got, want any) {
	t.Helper()

	switch want := want.(type) {
	case *EmptyLine:
		if _, ok := got.(*EmptyLine); !ok {
			t.Errorf("expected *EmptyLine, got %T", got)
		}
	case *TransactionHeaderLine:
		gotHeader, ok := got.(*TransactionHeaderLine)
		if !ok {
			t.Errorf("expected *TransactionHeaderLine, got %T", got)
			return
		}
		if gotHeader.Date != want.Date {
			t.Errorf("expected date %s, got %s", want.Date, gotHeader.Date)
		}
		if gotHeader.Description != want.Description {
			t.Errorf("expected description %s, got %s", want.Description, gotHeader.Description)
		}
	case *TransactionDetailLine:
		gotDetail, ok := got.(*TransactionDetailLine)
		if !ok {
			t.Errorf("expected *TransactionDetailLine, got %T", got)
			return
		}
		if gotDetail.Account != want.Account {
			t.Errorf("expected account %s, got %s", want.Account, gotDetail.Account)
		}
		if gotDetail.Amount != want.Amount {
			t.Errorf("expected amount %f, got %f", want.Amount, gotDetail.Amount)
		}
	default:
		t.Fatalf("unsupported want type %T", want)
	}
}
