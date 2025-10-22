package database

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestLogrusWriterPrintf(t *testing.T) {
	// Ensure Printf doesn't panic and logs via logrus
	var buf bytes.Buffer
	logrus.SetOutput(&buf)
	lw := newLogrusWriter(&buf)
	lw.Printf("hello %s", "world")
	if buf.Len() == 0 {
		t.Fatalf("expected output in buffer")
	}
}
