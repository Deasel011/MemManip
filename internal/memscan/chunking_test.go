package memscan

import "testing"

func TestFractureMemChunks_SplitsLargeRange(t *testing.T) {
	regions := []Region{{Base: 0x1000, Size: chunkSplitThreshold*2 + 7, Readable: true}}

	got := fractureMemChunks(regions)
	if len(got) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(got))
	}
	if got[0].Base != 0x1000 || got[0].Size != chunkSplitThreshold {
		t.Fatalf("first chunk mismatch: %+v", got[0])
	}
	if got[1].Base != 0x1000+chunkSplitThreshold || got[1].Size != chunkSplitThreshold {
		t.Fatalf("second chunk mismatch: %+v", got[1])
	}
	if got[2].Base != 0x1000+chunkSplitThreshold*2 || got[2].Size != 7 {
		t.Fatalf("tail chunk mismatch: %+v", got[2])
	}
}

func TestFractureMemChunks_LeavesSmallRangeUntouched(t *testing.T) {
	region := Region{Base: 0x2000, Size: 8, Readable: false}
	got := fractureMemChunks([]Region{region})
	if len(got) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(got))
	}
	if got[0] != region {
		t.Fatalf("region changed: got=%+v want=%+v", got[0], region)
	}
}
