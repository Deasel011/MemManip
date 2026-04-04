//go:build windows

package winapi

import "testing"

func TestFindProcessID_IntegrationStub(t *testing.T) {
	t.Skip("integration stub: requires a known running process on windows")
}

func TestOpenProcess_IntegrationStub(t *testing.T) {
	t.Skip("integration stub: requires a real pid and permissions on windows")
}

func TestListCommittedReadablePages_IntegrationStub(t *testing.T) {
	t.Skip("integration stub: requires open process handle on windows")
}
