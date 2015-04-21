package models

import (
	"regexp"
	"testing"
)

func TestMetaCast(t *testing.T) {
	var re = regexp.MustCompile("^(?:BINARY|CHAR|DATE|DATETIME|SIGNED|UNSIGNED|TIME|NUMERIC(?:\\(\\d+(?:,\\s?\\d+)?\\))?|DECIMAL(?:\\(\\d+(?:,\\s?\\d+)?\\))?)$")

	re.MatchString("NUMERIC(12)")
}
