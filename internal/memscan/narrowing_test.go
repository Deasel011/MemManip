package memscan

import (
	"fmt"
	"testing"
)

type fakeMemory struct {
	regions []Region
	reads   map[uintptr][]byte
	writes  map[uintptr][]byte
}

func (f *fakeMemory) OpenByPID(pid uint32) error { return nil }
func (f *fakeMemory) Close() error               { return nil }
func (f *fakeMemory) EnumerateRegions() ([]Region, error) {
	return f.regions, nil
}
func (f *fakeMemory) Read(address uintptr, size uintptr) ([]byte, error) {
	data, ok := f.reads[address]
	if !ok {
		return nil, fmt.Errorf("read miss at 0x%x", address)
	}
	if uintptr(len(data)) < size {
		return nil, fmt.Errorf("short read")
	}
	return data[:size], nil
}
func (f *fakeMemory) Write(address uintptr, data []byte) (int, error) {
	if f.writes == nil {
		f.writes = map[uintptr][]byte{}
	}
	f.writes[address] = append([]byte(nil), data...)
	return len(data), nil
}

func TestNarrow_FiltersExistingMatches(t *testing.T) {
	mem := &fakeMemory{reads: map[uintptr][]byte{
		0x1000: EncodeInt32LE(100),
		0x1004: EncodeInt32LE(200),
	}}
	s := NewScanner(mem)
	s.matches[0x1000] = struct{}{}
	s.matches[0x1004] = struct{}{}

	got, err := s.Narrow(EncodeInt32LE(200), 4)
	if err != nil {
		t.Fatalf("Narrow error: %v", err)
	}
	if len(got) != 1 || got[0] != 0x1004 {
		t.Fatalf("unexpected narrowed matches: %#v", got)
	}
}
