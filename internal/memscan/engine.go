package memscan

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
)

// Candidate is a matched memory address.
// Address is intentionally kept numeric and should only be formatted
// (for example as hex) at output/display boundaries.
type Candidate struct {
	Address uintptr
}

// ValueKind is a hook for future typed scan support.
type ValueKind int

const (
	KindInt16 ValueKind = iota
	KindInt32
	KindInt64
	KindFloat32
)

func valueSize(kind ValueKind) (uintptr, error) {
	switch kind {
	case KindInt16:
		return 2, nil
	case KindInt32:
		return 4, nil
	case KindInt64:
		return 8, nil
	case KindFloat32:
		return 4, nil
	default:
		return 0, fmt.Errorf("unsupported value kind %d", kind)
	}
}

func warn(op string, fields map[string]any) {
	msg := fmt.Sprintf("level=warn op=%s", op)
	for k, v := range fields {
		msg += fmt.Sprintf(" %s=%v", k, v)
	}
	log.Print(msg)
}

// SearchInt32 scans readable pages for target and returns candidates.
func (s *Scanner) SearchInt32(target int32) ([]Candidate, error) {
	regions, err := s.backend.EnumerateRegions()
	if err != nil {
		return nil, fmt.Errorf("enumerate regions: %w", err)
	}

	const op = "search_int32"
	regions = fractureMemChunks(regions)
	s.matches = make(map[uintptr]struct{})

	for _, region := range regions {
		if !region.Readable {
			warn(op, map[string]any{"reason": "region_not_readable", "base": region.Base, "size": region.Size})
			continue
		}
		if region.Size < 4 {
			continue
		}
		if region.Size > math.MaxInt {
			warn(op, map[string]any{"reason": "region_too_large", "base": region.Base, "size": region.Size})
			continue
		}

		data, err := s.backend.Read(region.Base, region.Size)
		if err != nil {
			warn(op, map[string]any{"reason": "read_failed", "base": region.Base, "size": region.Size, "error": err})
			continue
		}
		if uintptr(len(data)) != region.Size {
			warn(op, map[string]any{"reason": "partial_read", "base": region.Base, "expected": region.Size, "actual": len(data)})
			continue
		}

		for off := 0; off+4 <= len(data); off += 4 {
			v := int32(binary.LittleEndian.Uint32(data[off : off+4]))
			if v == target {
				s.matches[region.Base+uintptr(off)] = struct{}{}
			}
		}
	}

	return s.candidatesFromMatches(), nil
}

// NarrowInt32 filters existing candidates to addresses still equal to target.
func (s *Scanner) NarrowInt32(candidates []Candidate, target int32) ([]Candidate, error) {
	if len(candidates) == 0 {
		return nil, errors.New("no candidates provided")
	}

	const op = "narrow_int32"
	next := make(map[uintptr]struct{}, len(candidates))
	for _, c := range candidates {
		data, err := s.backend.Read(c.Address, 4)
		if err != nil {
			warn(op, map[string]any{"reason": "read_failed", "address": c.Address, "error": err})
			continue
		}
		if len(data) != 4 {
			warn(op, map[string]any{"reason": "partial_read", "address": c.Address, "expected": 4, "actual": len(data)})
			continue
		}
		if int32(binary.LittleEndian.Uint32(data)) == target {
			next[c.Address] = struct{}{}
		}
	}

	s.matches = next
	return s.candidatesFromMatches(), nil
}

// SetInt32 writes value to every candidate address.
func (s *Scanner) SetInt32(candidates []Candidate, value int32) (int, error) {
	if len(candidates) == 0 {
		return 0, errors.New("no candidates provided")
	}

	const op = "set_int32"
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(value))

	updated := 0
	for _, c := range candidates {
		n, err := s.backend.Write(c.Address, buf)
		if err != nil {
			warn(op, map[string]any{"reason": "write_failed", "address": c.Address, "error": err})
			continue
		}
		if n != len(buf) {
			warn(op, map[string]any{"reason": "partial_write", "address": c.Address, "expected": len(buf), "actual": n})
			continue
		}
		updated++
	}

	return updated, nil
}

func (s *Scanner) candidatesFromMatches() []Candidate {
	matches := s.Matches()
	out := make([]Candidate, 0, len(matches))
	for _, addr := range matches {
		out = append(out, Candidate{Address: addr})
	}
	return out
}
