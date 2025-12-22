package params

import (
	"strings"

	derrors "github.com/KlimKlimKlimKlim/trossage-backend/internal/errors"
)

const (
	MinSearchQueryLength = 3
)

func ParseSearchQuery(s string) (string, error) {
	query := strings.ToLower(strings.TrimSpace(s))

	if query == "" {
		return "", derrors.ErrEmptyQuery
	}

	if len(query) < MinSearchQueryLength {
		return "", derrors.ErrSearchQueryTooShort
	}

	return query, nil
}
