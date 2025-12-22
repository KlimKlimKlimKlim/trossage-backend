package params

import (
	"strconv"
)

const (
	DefaultLimit = 20
	MinLimit     = 1
	MaxLimit     = 100

	DefaultOffset = 0
)

func ParseLimit(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil || value < MinLimit || value > MaxLimit {
		return DefaultLimit
	}

	return value
}

func ParseOffset(s string) int {
	value, err := strconv.Atoi(s)
	if err != nil || value < DefaultOffset {
		return DefaultOffset
	}

	return value
}
