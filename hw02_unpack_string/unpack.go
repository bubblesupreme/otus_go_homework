package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"

	log "github.com/sirupsen/logrus"
)

var ErrInvalidString = errors.New("invalid string")
var ErrWriteSting = errors.New("invalid string")

func writeStrToBuilder(builder *strings.Builder, str string) error {
	if _, err := builder.WriteString(str); err != nil {
		log.WithField("string to write", str).Error(err)
		return ErrWriteSting
	}

	return nil
}

// Write substring of str to builder.
func flushSubStr(str string, builder *strings.Builder, firstSubIdx int, lastSubIdx int) error {
	return writeStrToBuilder(builder, str[firstSubIdx:lastSubIdx])
}

// Repeat symbol {times} times and write.
func writeRepeat(symbol rune, builder *strings.Builder, times int) error {
	return writeStrToBuilder(builder, strings.Repeat(string(symbol), times))
}

// Read symbols before find a digit and keep the index of the last written symbol
// then write all symbols started from that index, multiplying the last symbol by digit.
func Unpack(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	if unicode.IsDigit(rune(str[0])) || str[0] == '\\' {
		return "", ErrInvalidString
	}

	builder := strings.Builder{}
	escaping := false
	firstSubIdx := 0
	for i := 1; i < len(str); i++ {
		ch := rune(str[i])

		if escaping {
			if !unicode.IsDigit(ch) && ch != '\\' {
				return "", ErrInvalidString
			}

			if err := flushSubStr(str, &builder, firstSubIdx, i-1); err != nil {
				return "", err
			}
			firstSubIdx = i
			escaping = false
			continue
		}

		if unicode.IsDigit(ch) {
			// if 2 unescaped digits in a row
			if i == firstSubIdx {
				return "", ErrInvalidString
			}

			if err := flushSubStr(str, &builder, firstSubIdx, i-1); err != nil {
				return "", err
			}

			if err := writeRepeat(rune(str[i-1]), &builder, int(str[i]-'0')); err != nil {
				return "", err
			}
			firstSubIdx = i + 1
			continue
		}

		if ch == '\\' {
			escaping = true
		}
	}

	if escaping {
		return "", ErrInvalidString
	}

	if err := flushSubStr(str, &builder, firstSubIdx, len(str)); err != nil {
		return "", err
	}

	return builder.String(), nil
}
