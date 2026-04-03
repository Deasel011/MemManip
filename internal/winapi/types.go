package winapi

// Handle mirrors a Windows HANDLE.
type Handle uintptr

// MemoryBasicInformation matches the Win32 MEMORY_BASIC_INFORMATION layout.
type MemoryBasicInformation struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	PartitionID       uint16
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

// SystemInfo matches the subset of Win32 SYSTEM_INFO used by memory scanning.
type SystemInfo struct {
	ProcessorArchitecture uint16
	Reserved              uint16
	PageSize              uint32
	MinimumAddress        uintptr
	MaximumAddress        uintptr
	ActiveProcessorMask   uintptr
	NumberOfProcessors    uint32
	ProcessorType         uint32
	AllocationGranularity uint32
	ProcessorLevel        uint16
	ProcessorRevision     uint16
}
