package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"

	"memmanip/internal/memscan"
	"memmanip/internal/winapi"
)

func main() {
	var (
		pid     = flag.Uint("pid", 0, "target process ID")
		mode    = flag.String("mode", "search", "operation: search|narrow|set")
		value   = flag.Int("value", 0, "target 32-bit integer value")
		initial = flag.Int("initial", 0, "initial value used before narrow/set")
		stride  = flag.Int("stride", 4, "scan step in bytes")
	)
	flag.Parse()

	if *pid == 0 {
		log.Fatal("-pid is required")
	}

	backend := winapi.NewProcess()
	scanner := memscan.NewScanner(backend)
	if err := scanner.Open(uint32(*pid)); err != nil {
		log.Fatalf("open process: %v", err)
	}
	defer scanner.Close()

	switch *mode {
	case "search":
		matches, err := scanner.Search(int32LE(*value), *stride)
		if err != nil {
			log.Fatalf("search failed: %v", err)
		}
		printMatches(matches)
	case "narrow":
		if _, err := scanner.Search(int32LE(*initial), *stride); err != nil {
			log.Fatalf("seed search failed: %v", err)
		}
		matches, err := scanner.Narrow(int32LE(*value), *stride)
		if err != nil {
			log.Fatalf("narrow failed: %v", err)
		}
		printMatches(matches)
	case "set":
		if _, err := scanner.Search(int32LE(*initial), *stride); err != nil {
			log.Fatalf("seed search failed: %v", err)
		}
		updated, err := scanner.Set(int32LE(*value))
		if err != nil {
			log.Fatalf("set failed: %v", err)
		}
		fmt.Printf("updated %d addresses\n", updated)
	default:
		log.Fatalf("unsupported mode %q", *mode)
	}
}

func int32LE(v int) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(int32(v)))
	return buf
}

func printMatches(matches []uintptr) {
	fmt.Printf("matches: %d\n", len(matches))
	for _, addr := range matches {
		fmt.Printf("0x%X\n", addr)
	}
}
