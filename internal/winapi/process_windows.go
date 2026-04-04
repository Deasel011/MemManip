//go:build windows

package winapi

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/sys/windows"
)

const defaultProcessAccess = ProcessQueryInformation | ProcessVMRead | ProcessVMWrite | ProcessVMOperation

// FindProcessID locates a process by executable name using the Toolhelp snapshot APIs.
func FindProcessID(exeName string) (uint32, error) {
	trimmedName := strings.TrimSpace(exeName)
	if trimmedName == "" {
		return 0, fmt.Errorf("FindProcessID: executable name is empty")
	}

	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, fmt.Errorf("CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS): %w (win32=%d)", err, extractErrno(err))
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(windows.SizeofProcessEntry32)

	if err := windows.Process32First(snapshot, &entry); err != nil {
		return 0, fmt.Errorf("Process32First: %w (win32=%d)", err, extractErrno(err))
	}

	target := strings.ToLower(trimmedName)
	for {
		name := windows.UTF16ToString(entry.ExeFile[:])
		if strings.EqualFold(name, target) {
			return entry.ProcessID, nil
		}

		err = windows.Process32Next(snapshot, &entry)
		if err == nil {
			continue
		}
		if errors.Is(err, windows.ERROR_NO_MORE_FILES) {
			return 0, fmt.Errorf("FindProcessID(%q): process not found", exeName)
		}
		return 0, fmt.Errorf("Process32Next: %w (win32=%d)", err, extractErrno(err))
	}
}

// OpenProcess opens a process with requested access rights.
func OpenProcess(pid uint32, access uint32) (windows.Handle, error) {
	if pid == 0 {
		return 0, fmt.Errorf("OpenProcess: pid must be non-zero")
	}

	requiredAccess := access
	if requiredAccess == 0 {
		requiredAccess = defaultProcessAccess
	}

	handle, err := windows.OpenProcess(requiredAccess, false, pid)
	if err != nil {
		return 0, fmt.Errorf("OpenProcess(pid=%d, access=0x%x): %w (win32=%d)", pid, requiredAccess, err, extractErrno(err))
	}

	closeOnError := true
	defer func() {
		if closeOnError {
			_ = windows.CloseHandle(handle)
		}
	}()

	closeOnError = false
	return handle, nil
}

func extractErrno(err error) uint32 {
	var errno windows.Errno
	if errors.As(err, &errno) {
		return uint32(errno)
	}
	return 0
}
