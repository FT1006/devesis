package main

import (
	"bytes"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run main
	main()

	// Restore stdout
	w.Close()
	os.Stdout = old

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	expected := "Board initialised: 6×7 rooms\n"
	if output != expected {
		t.Errorf("Expected %q, got %q", expected, output)
	}
}