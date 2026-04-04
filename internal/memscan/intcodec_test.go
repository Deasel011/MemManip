package memscan

import "testing"

func TestEncodeDecodeInt32LE_RoundTrip(t *testing.T) {
	in := int32(-1234567)
	encoded := EncodeInt32LE(in)
	got, err := DecodeInt32LE(encoded)
	if err != nil {
		t.Fatalf("DecodeInt32LE returned error: %v", err)
	}
	if got != in {
		t.Fatalf("round trip mismatch: got=%d want=%d", got, in)
	}
}

func TestDecodeInt32LE_InvalidSize(t *testing.T) {
	if _, err := DecodeInt32LE([]byte{1, 2, 3}); err == nil {
		t.Fatal("expected error for non-4-byte input")
	}
}
