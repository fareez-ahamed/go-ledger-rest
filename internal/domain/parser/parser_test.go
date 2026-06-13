package parser

import (
	"fmt"
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
