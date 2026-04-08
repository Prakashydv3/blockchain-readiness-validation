package main

import "fmt"

// ReplaySystem verifies envelope integrity, recomputes state root, and validates anchor.
// Returns (true, stateRoot) on success, (false, reason) on any failure.
func ReplaySystem(envelopes []Envelope, anchor Anchor) (bool, string) {
	// 1. Verify every envelope hash
	for _, e := range envelopes {
		if hashEnvelope(e) != e.Hash {
			return false, fmt.Sprintf("envelope hash mismatch: %s", e.ExecutionID)
		}
	}

	// 2. Recompute and verify state root
	stateRoot := GenerateStateRoot(envelopes)
	if stateRoot != anchor.StateRoot {
		return false, fmt.Sprintf("state root mismatch: got %s want %s", stateRoot[:12], anchor.StateRoot[:12])
	}

	// 3. Verify anchor hash
	if hashAnchor(anchor.StateRoot, anchor.Timestamp) != anchor.Hash {
		return false, "anchor hash mismatch"
	}

	return true, stateRoot
}
