package main

import (
	"flag"
	"fmt"
	"log"

	"memmanip/internal/memscan"
	"memmanip/internal/winapi"
)

type optionalInt struct {
	set   bool
	value int
}

func (o *optionalInt) Set(s string) error {
	o.set = true
	_, err := fmt.Sscanf(s, "%d", &o.value)
	if err != nil {
		return fmt.Errorf("invalid integer %q", s)
	}
	return nil
}

func (o *optionalInt) String() string {
	if !o.set {
		return ""
	}
	return fmt.Sprintf("%d", o.value)
}

func main() {
	var (
		processName string
		searchValue int
		narrowValue optionalInt
		setValue    optionalInt
		size        int
	)

	flag.StringVar(&processName, "process", "", "target process executable name (example: Darkest.exe)")
	flag.IntVar(&searchValue, "search", 0, "initial int value to search for")
	flag.Var(&narrowValue, "narrow", "optional int value to narrow existing results")
	flag.Var(&setValue, "set", "optional int value to write at remaining addresses")
	flag.IntVar(&size, "size", 4, "value size in bytes")
	flag.Parse()

	if processName == "" {
		log.Fatal("-process is required")
	}
	if size != 4 {
		log.Fatalf("unsupported -size=%d (only 4-byte int32 is currently supported)", size)
	}

	pid, err := winapi.FindProcessID(processName)
	if err != nil {
		log.Fatalf("find process id: %v", err)
	}
	fmt.Printf("process=%s pid=%d\n", processName, pid)

	backend := winapi.NewProcess()
	scanner := memscan.NewScanner(backend)
	if err := scanner.Open(pid); err != nil {
		log.Fatalf("open process: %v", err)
	}
	defer scanner.Close()

	matches, err := scanner.Search(memscan.EncodeInt32LE(int32(searchValue)), size)
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}
	printMatches("search", matches)

	if narrowValue.set {
		matches, err = scanner.Narrow(memscan.EncodeInt32LE(int32(narrowValue.value)), size)
		if err != nil {
			log.Fatalf("narrow failed: %v", err)
		}
		printMatches("narrow", matches)
	}

	if setValue.set {
		updated, err := scanner.Set(memscan.EncodeInt32LE(int32(setValue.value)))
		if err != nil {
			log.Fatalf("set failed: %v", err)
		}
		fmt.Printf("set updated %d addresses\n", updated)
	}
}

func printMatches(step string, matches []uintptr) {
	fmt.Printf("%s matches: %d\n", step, len(matches))
	for _, addr := range matches {
		fmt.Printf("0x%08X\n", uint32(addr))
	}
}
