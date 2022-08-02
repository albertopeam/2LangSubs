package main

import (
	"os"
	"testing"

	"github.com/udhos/equalfile"
)

func TestGivenInputArgsWhenInvokeMainThenProcessSubtitlesAndMatchOutputWithExpectedOne(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-s1=testdata/test_input_s1.srt", "-s2=testdata/test_input_s2.srt", "-sO=testdata/test_output.srt"}

	main()

	cmp := equalfile.New(nil, equalfile.Options{})
	equal, err := cmp.CompareFile("testdata/test_output.srt", "testdata/test_output_s1+s2.srt")

	if !equal {
		t.Fatalf("Generated output and expected output are not equals, %v", err)
	}

	path := "./testdata/test_output.srt"
	e := os.Remove(path)
	if e != nil {
		t.Fatalf("Remove %v throws an error %v", path, e)
	}
}
