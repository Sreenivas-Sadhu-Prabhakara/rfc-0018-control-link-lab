package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"example.com/rfc-go-labs/lab"
)

func TestSampleOutputs(t *testing.T) {
	input, err := os.ReadFile("../../examples/burst.json")
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := run(nil, bytes.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	var got lab.Report
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Classes) != 2 || len(got.Results) != 9 ||
		got.Classes[0] != (lab.Summary{Class: lab.Host, Count: 3, SharedMax: 3, SplitMax: 1}) {
		t.Fatalf("unexpected summary: %+v", got)
	}
	out.Reset()
	if err := run([]string{"-html"}, bytes.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "<th scope=\"row\">host</th><td>3</td><td>3</td><td>1</td>") {
		t.Fatal("HTML report is missing expected host summary")
	}
}

func TestCLIRejectsMalformedInput(t *testing.T) {
	for _, input := range []string{
		`{"arrivals":[{"tick":0,"class":"host","extra":true}]}`,
		`{"arrivals":[{"tick":0,"class":"host"}]} {}`,
		strings.Repeat(" ", (1<<20)+1),
	} {
		var out bytes.Buffer
		if err := run(nil, strings.NewReader(input), &out); err == nil || out.Len() != 0 {
			t.Errorf("invalid input of %d bytes gave %d output bytes and error %v", len(input), out.Len(), err)
		}
	}
}
