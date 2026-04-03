package memscan

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

// Region is a platform-neutral view of a process memory segment.
type Region struct {
	Base     uintptr
	Size     uintptr
	Readable bool
}

// ProcessMemory is the minimal OS-backed contract required by Scanner.
type ProcessMemory interface {
	OpenByPID(pid uint32) error
	Close() error
	EnumerateRegions() ([]Region, error)
	Read(address uintptr, size uintptr) ([]byte, error)
	Write(address uintptr, data []byte) (int, error)
}

// Scanner owns stateful match addresses so Search, Narrow, and Set can be chained.
type Scanner struct {
	backend ProcessMemory
	matches map[uintptr]struct{}
}

func NewScanner(backend ProcessMemory) *Scanner {
	return &Scanner{backend: backend, matches: make(map[uintptr]struct{})}
}

func (s *Scanner) Open(pid uint32) error {
	return s.backend.OpenByPID(pid)
}

func (s *Scanner) Close() error {
	return s.backend.Close()
}

func (s *Scanner) Search(value []byte, stride int) ([]uintptr, error) {
	if len(value) == 0 {
		return nil, errors.New("value must not be empty")
	}
	if stride <= 0 {
		return nil, errors.New("stride must be > 0")
	}

	regions, err := s.backend.EnumerateRegions()
	if err != nil {
		return nil, fmt.Errorf("enumerate regions: %w", err)
	}

	s.matches = make(map[uintptr]struct{})
	for _, region := range regions {
		if !region.Readable || region.Size < uintptr(len(value)) {
			continue
		}

		data, err := s.backend.Read(region.Base, region.Size)
		if err != nil {
			continue
		}

		for off := 0; off+len(value) <= len(data); off += stride {
			if bytes.Equal(data[off:off+len(value)], value) {
				s.matches[region.Base+uintptr(off)] = struct{}{}
			}
		}
	}

	return s.Matches(), nil
}

func (s *Scanner) Narrow(value []byte, stride int) ([]uintptr, error) {
	if len(s.matches) == 0 {
		return nil, errors.New("no prior matches, run Search first")
	}
	if len(value) == 0 {
		return nil, errors.New("value must not be empty")
	}
	if stride <= 0 {
		return nil, errors.New("stride must be > 0")
	}

	next := make(map[uintptr]struct{})
	for addr := range s.matches {
		data, err := s.backend.Read(addr, uintptr(len(value)))
		if err != nil {
			continue
		}
		if bytes.Equal(data, value) {
			next[addr] = struct{}{}
		}
	}

	s.matches = next
	return s.Matches(), nil
}

func (s *Scanner) Set(value []byte) (int, error) {
	if len(s.matches) == 0 {
		return 0, errors.New("no prior matches, run Search first")
	}
	if len(value) == 0 {
		return 0, errors.New("value must not be empty")
	}

	updated := 0
	for addr := range s.matches {
		n, err := s.backend.Write(addr, value)
		if err != nil {
			continue
		}
		if n == len(value) {
			updated++
		}
	}
	return updated, nil
}

func (s *Scanner) Matches() []uintptr {
	out := make([]uintptr, 0, len(s.matches))
	for addr := range s.matches {
		out = append(out, addr)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
