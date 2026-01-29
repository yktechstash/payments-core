package domain

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Money struct {
	decimal.Decimal
}

func ParseMoney(s string) (Money, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Money{}, fmt.Errorf("%w: amount is empty", ErrInvalidInput)
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Money{}, fmt.Errorf("%w: invalid decimal", ErrInvalidInput)
	}
	return Money{Decimal: d}, nil
}

func (m Money) IsNegativeOrZero() bool {
	return m.Cmp(decimal.Zero) <= 0
}

func (m Money) String() string {
	return m.Decimal.String()
}
