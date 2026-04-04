package memscan

import "math"

const chunkSplitThreshold uintptr = 1921024

// fractureMemChunks splits large regions into stable, non-overlapping chunks
// while preserving original ordering.
func fractureMemChunks(regions []Region) []Region {
	if len(regions) == 0 {
		return nil
	}

	out := make([]Region, 0, len(regions))
	for _, region := range regions {
		if region.Size == 0 || region.Size <= chunkSplitThreshold {
			out = append(out, region)
			continue
		}

		base := region.Base
		remaining := region.Size
		maxParts := int((uint64(region.Size) / uint64(chunkSplitThreshold)) + 1)
		if maxParts <= 0 {
			maxParts = 1
		}

		for parts := 0; remaining > 0 && parts < maxParts; parts++ {
			chunkSize := remaining
			if chunkSize > chunkSplitThreshold {
				chunkSize = chunkSplitThreshold
			}

			nextBase, ok := addUintptr(base, chunkSize)
			if !ok {
				// Clamp to the maximum representable range to avoid overflow loops.
				chunkSize = math.MaxUint - base
				if chunkSize == 0 {
					break
				}
				nextBase = base + chunkSize
			}

			out = append(out, Region{
				Base:     base,
				Size:     chunkSize,
				Readable: region.Readable,
			})

			base = nextBase
			remaining -= chunkSize
		}
	}

	return out
}

func addUintptr(a, b uintptr) (uintptr, bool) {
	sum := a + b
	return sum, sum >= a
}
