//go:build !windows

package winapi

import "fmt"

func FindProcessID(exeName string) (uint32, error) {
	return 0, fmt.Errorf("FindProcessID is only available on windows")
}
