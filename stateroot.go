package main

import (
	"sort"
)

// GenerateStateRoot sorts envelopes by ExecutionID then hashes all envelope hashes concatenated.
// SAME input (any order) → SAME root always.
func GenerateStateRoot(envelopes []Envelope) string {
	sorted := make([]Envelope, len(envelopes))
	copy(sorted, envelopes)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ExecutionID < sorted[j].ExecutionID
	})

	combined := ""
	for _, e := range sorted {
		combined += e.Hash
	}
	return sha256Hex(combined)
}
