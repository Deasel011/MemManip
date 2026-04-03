//go:build windows

package winapi

import (
	"fmt"
	"syscall"
	"unsafe"

	"memmanip/internal/memscan"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess        = kernel32.NewProc("OpenProcess")
	procCloseHandle        = kernel32.NewProc("CloseHandle")
	procReadProcessMemory  = kernel32.NewProc("ReadProcessMemory")
	procWriteProcessMemory = kernel32.NewProc("WriteProcessMemory")
	procVirtualQueryEx     = kernel32.NewProc("VirtualQueryEx")
	procGetSystemInfo      = kernel32.NewProc("GetSystemInfo")
)

type Process struct {
	handle Handle
}

func NewProcess() *Process {
	return &Process{}
}

func (p *Process) OpenByPID(pid uint32) error {
	r1, _, err := procOpenProcess.Call(
		uintptr(ProcessQueryInformation|ProcessVMRead|ProcessVMWrite|ProcessVMOperation),
		0,
		uintptr(pid),
	)
	if r1 == 0 {
		return fmt.Errorf("OpenProcess(%d): %w", pid, err)
	}
	p.handle = Handle(r1)
	return nil
}

func (p *Process) Close() error {
	if p.handle == 0 {
		return nil
	}
	r1, _, err := procCloseHandle.Call(uintptr(p.handle))
	if r1 == 0 {
		return fmt.Errorf("CloseHandle: %w", err)
	}
	p.handle = 0
	return nil
}

func (p *Process) EnumerateRegions() ([]memscan.Region, error) {
	if p.handle == 0 {
		return nil, fmt.Errorf("process is not open")
	}

	var si SystemInfo
	procGetSystemInfo.Call(uintptr(unsafe.Pointer(&si)))

	regions := make([]memscan.Region, 0, 1024)
	for addr := si.MinimumAddress; addr < si.MaximumAddress; {
		var mbi MemoryBasicInformation
		r1, _, _ := procVirtualQueryEx.Call(
			uintptr(p.handle),
			addr,
			uintptr(unsafe.Pointer(&mbi)),
			unsafe.Sizeof(mbi),
		)
		if r1 == 0 {
			break
		}

		regions = append(regions, memscan.Region{
			Base:     mbi.BaseAddress,
			Size:     mbi.RegionSize,
			Readable: isReadable(mbi),
		})

		next := mbi.BaseAddress + mbi.RegionSize
		if next <= addr {
			break
		}
		addr = next
	}

	return regions, nil
}

func (p *Process) Read(address uintptr, size uintptr) ([]byte, error) {
	if p.handle == 0 {
		return nil, fmt.Errorf("process is not open")
	}
	if size == 0 {
		return []byte{}, nil
	}

	buf := make([]byte, size)
	var read uintptr
	r1, _, err := procReadProcessMemory.Call(
		uintptr(p.handle),
		address,
		uintptr(unsafe.Pointer(&buf[0])),
		size,
		uintptr(unsafe.Pointer(&read)),
	)
	if r1 == 0 {
		return nil, fmt.Errorf("ReadProcessMemory(0x%x): %w", address, err)
	}
	return buf[:read], nil
}

func (p *Process) Write(address uintptr, data []byte) (int, error) {
	if p.handle == 0 {
		return 0, fmt.Errorf("process is not open")
	}
	if len(data) == 0 {
		return 0, nil
	}

	var written uintptr
	r1, _, err := procWriteProcessMemory.Call(
		uintptr(p.handle),
		address,
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
	)
	if r1 == 0 {
		return int(written), fmt.Errorf("WriteProcessMemory(0x%x): %w", address, err)
	}
	return int(written), nil
}

func isReadable(mbi MemoryBasicInformation) bool {
	if mbi.State != MemCommit {
		return false
	}
	if mbi.Protect&PageGuard != 0 || mbi.Protect&PageNoAccess != 0 {
		return false
	}
	return mbi.Protect&readableProtectMask != 0
}
