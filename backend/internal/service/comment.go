package service

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxCommentLength = 500

// cleanComment cleans and validates a donor's comment (empty = no comment).
func cleanComment(comment string) (string, error) {
	comment = strings.TrimSpace(comment)
	if comment == "" {
		return "", nil
	}
	if utf8.RuneCountInString(comment) > maxCommentLength {
		return "", errors.New("comment is too long (500 characters maximum)")
	}
	for _, r := range comment {
		if unicode.IsControl(r) && r != '\n' && r != '\r' {
			return "", errors.New("comment contains invalid characters")
		}
	}

	lower := strings.ToLower(comment)
	if strings.Contains(lower, "http://") || strings.Contains(lower, "https://") || strings.Contains(lower, "www.") {
		return "", errors.New("links are not allowed in comments")
	}

	return comment, nil
}
