package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

var testOk = `1
2
3
3
4
5
`

var testOkResult = `1
2
3
4
5
`

func TestOk(t *testing.T) {

	in := bufio.NewReader(strings.NewReader(testOk))
	out := new(bytes.Buffer)
	err := uniq(in, out)
	if err != nil {
		t.Errorf("test for Ok Failed - error")
	}

	result := out.String()
	if result != testOkResult {
		t.Errorf("test for Ok Failed - results not matched\n %v %v", result, testOkResult)
	}
}

var testFailed = `1
2
1
`

func TestForError(t *testing.T) {
	in := bufio.NewReader(strings.NewReader(testFailed))
	out := new(bytes.Buffer)
	err := uniq(in, out)
	if err == nil {
		t.Errorf("test for Ok Failed - error: %v", err)
	}
}
