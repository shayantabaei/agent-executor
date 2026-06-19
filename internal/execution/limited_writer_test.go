package execution

import "testing"

func TestLimitedWriterTruncatesOutput(t *testing.T) {
	writer := NewLimitedWriter(10)

	n, err := writer.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Fatalf("expected 5 bytes written, got %d", n)
	}

	n, err = writer.Write([]byte(" world"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 6 {
		t.Fatalf("expected writer to report 6 bytes consumed, got %d", n)
	}

	got := writer.String()
	want := "hello worl"

	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestLimitedWriterIgnoresWritesAfterLimit(t *testing.T) {
	writer := NewLimitedWriter(5)

	_, _ = writer.Write([]byte("hello"))
	n, err := writer.Write([]byte("world"))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 5 {
		t.Fatalf("expected 5 bytes consumed, got %d", n)
	}
	if got := writer.String(); got != "hello" {
		t.Fatalf("expected %q, got %q", "hello", got)
	}
}
