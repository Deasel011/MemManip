//go:build windows

package winapi

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Page represents one committed, readable page range in a remote process.
type Page struct {
	BaseAddress uintptr
	Size        uintptr
	Protect     uint32
}

// ListCommittedReadablePages enumerates a process address space via VirtualQueryEx
// and returns only committed readable ranges in ascending address order.
func ListCommittedReadablePages(handle windows.Handle) ([]Page, error) {
	if handle == 0 {
		return nil, fmt.Errorf("process handle is invalid")
	}

	var si SystemInfo
	procGetSystemInfo.Call(uintptr(unsafe.Pointer(&si)))

	pages := make([]Page, 0, 1024)
	for addr := si.MinimumAddress; addr < si.MaximumAddress; {
		var mbi MemoryBasicInformation
		r1, _, err := procVirtualQueryEx.Call(
			uintptr(handle),
			addr,
			uintptr(unsafe.Pointer(&mbi)),
			unsafe.Sizeof(mbi),
		)
		if r1 == 0 {
			if len(pages) == 0 {
				return nil, fmt.Errorf("VirtualQueryEx(0x%x): %w", addr, err)
			}
			break
		}

		if isReadable(mbi) {
			pages = append(pages, Page{
				BaseAddress: mbi.BaseAddress,
				Size:        mbi.RegionSize,
				Protect:     mbi.Protect,
			})
		}

		if mbi.RegionSize == 0 {
			break
		}
		next := mbi.BaseAddress + mbi.RegionSize
		if next <= addr {
			break
		}
		addr = next
	}

	return pages, nil
}
