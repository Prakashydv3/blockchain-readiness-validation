package main

import (
	"fmt"
	"time"
)

// ExecutionAgent picks an execution and produces a sealed Envelope.
func ExecutionAgent(exec Execution) Envelope {
	inputHash := sha256Hex(exec.Input)
	outputHash := sha256Hex(exec.Output)
	env := Envelope{
		ExecutionID: exec.ID,
		InputHash:   inputHash,
		OutputHash:  outputHash,
	}
	env.Hash = hashEnvelope(env)
	fmt.Printf("[ExecutionAgent] sealed  id=%s  hash=%s\n", env.ExecutionID, env.Hash[:12])
	return env
}

// ValidationAgent replays an envelope against the original execution.
// Returns false + reason if any hash mismatches.
func ValidationAgent(env Envelope, exec Execution) (bool, string) {
	if sha256Hex(exec.Input) != env.InputHash {
		return false, fmt.Sprintf("input hash mismatch for %s", env.ExecutionID)
	}
	if sha256Hex(exec.Output) != env.OutputHash {
		return false, fmt.Sprintf("output hash mismatch for %s", env.ExecutionID)
	}
	if hashEnvelope(env) != env.Hash {
		return false, fmt.Sprintf("envelope hash mismatch for %s", env.ExecutionID)
	}
	fmt.Printf("[ValidationAgent] valid  id=%s\n", env.ExecutionID)
	return true, ""
}

// RelayAgent collects validated envelopes, generates state root, and anchors to L1.
func RelayAgent(envelopes []Envelope) Anchor {
	stateRoot := GenerateStateRoot(envelopes)
	ts := time.Now().UnixNano()
	anchor := Anchor{
		StateRoot: stateRoot,
		Timestamp: ts,
		Hash:      hashAnchor(stateRoot, ts),
	}
	fmt.Printf("[RelayAgent]  state_root=%s\n", stateRoot[:12])
	fmt.Printf("[RelayAgent]  anchor_hash=%s\n", anchor.Hash[:12])
	return anchor
}
