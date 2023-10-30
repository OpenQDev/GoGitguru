package gitutil

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func TestSplitWithDelimiters(t *testing.T) {
	text := "commit 1234567890abcdef"
	pattern := regexp.MustCompile(`commit [a-f0-9]{40}`)
	expected := []string{"commit 1234567890abcdef"}

	result := SplitWithDelimiters(text, pattern)
	fmt.Println(result)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
