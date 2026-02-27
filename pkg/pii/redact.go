// Package pii provides personally identifiable information (PII) handling functionality.
package pii

import (
	"strings"
)

var redactReplacement = []byte("***")

// Redact - masks sensitive string by replacing some characters in the middle with ***.
func Redact(src string) string {
	if src == "" {
		return ""
	}

	newLen := len(redactReplacement) + 2

	index := strings.IndexByte(src, '@')
	if index > 0 {
		newLen += len(src) - index
	}

	var res []byte
	// reuse the src buffer if it's big enough so that we don't have to allocate
	if newLen <= len(src) {
		res = []byte(src)
	} else {
		res = make([]byte, newLen)
	}
	res[0] = src[0]
	copy(res[1:], redactReplacement)

	if index > 0 {
		res[len(redactReplacement)+1] = src[index-1]
		copy(res[len(redactReplacement)+2:], []byte(src[index:]))
	} else {
		res[len(redactReplacement)+1] = src[len(src)-1]
	}

	return string(res[0:newLen])
}

// RedactAll masks all strings by replacing some characters with ***.
func RedactAll(src []string) []string {
	res := make([]string, len(src))
	for i := range src {
		res[i] = Redact(src[i])
	}

	return res
}
