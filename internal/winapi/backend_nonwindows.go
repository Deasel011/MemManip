//go:build !windows

package winapi

import (
	"fmt"

	"memmanip/internal/memscan"
)

type Process struct{}

func NewProcess() *Process { return &Process{} }

func (p *Process) OpenByPID(pid uint32) error {
	return fmt.Errorf("winapi backend is only available on windows")
}

func (p *Process) Close() error { return nil }

func (p *Process) EnumerateRegions() ([]memscan.Region, error) {
	return nil, fmt.Errorf("winapi backend is only available on windows")
}

func (p *Process) Read(address uintptr, size uintptr) ([]byte, error) {
	return nil, fmt.Errorf("winapi backend is only available on windows")
}

func (p *Process) Write(address uintptr, data []byte) (int, error) {
	return 0, fmt.Errorf("winapi backend is only available on windows")
}
