package winapi

const (
	ProcessQueryInformation = 0x0400
	ProcessVMRead           = 0x0010
	ProcessVMWrite          = 0x0020
	ProcessVMOperation      = 0x0008

	MemCommit  = 0x00001000
	MemReserve = 0x00002000

	PageNoAccess         = 0x01
	PageReadOnly         = 0x02
	PageReadWrite        = 0x04
	PageWriteCopy        = 0x08
	PageExecute          = 0x10
	PageExecuteRead      = 0x20
	PageExecuteReadWrite = 0x40
	PageExecuteWriteCopy = 0x80
	PageGuard            = 0x100
	PageNoCache          = 0x200
	PageWriteCombine     = 0x400
)

const readableProtectMask = PageReadOnly | PageReadWrite | PageWriteCopy | PageExecuteRead | PageExecuteReadWrite | PageExecuteWriteCopy
