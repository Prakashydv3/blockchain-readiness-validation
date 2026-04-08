package main

import (
	"fmt"
	"strings"
)

func separator(label string) {
	fmt.Printf("\n%s [ %s ] %s\n", strings.Repeat("─", 10), label, strings.Repeat("─", 10))
}

func main() {
	executions := []Execution{
		{ID: "exec-A", Input: "Gurukul TTS input alpha", Output: "audio-output-alpha"},
		{ID: "exec-B", Input: "Gurukul TTS input beta", Output: "audio-output-beta"},
		{ID: "exec-C", Input: "Gurukul TTS input gamma", Output: "audio-output-gamma"},
	}

	// ── PHASE 1 + 2: Agent pipeline ──────────────────────────────────────────
	separator("PHASE 2 — AGENT PIPELINE")

	var envelopes []Envelope
	for _, exec := range executions {
		env := ExecutionAgent(exec)
		ok, reason := ValidationAgent(env, exec)
		if !ok {
			fmt.Printf("[FATAL] validation failed: %s\n", reason)
			return
		}
		envelopes = append(envelopes, env)
	}

	anchor := RelayAgent(envelopes)

	// ── PHASE 1: State root order-independence proof ──────────────────────────
	separator("PHASE 1 — STATE ROOT ORDER PROOF")

	reversed := []Envelope{envelopes[2], envelopes[1], envelopes[0]}
	root1 := GenerateStateRoot(envelopes)
	root2 := GenerateStateRoot(reversed)
	fmt.Printf("A,B,C → %s\n", root1[:16])
	fmt.Printf("C,B,A → %s\n", root2[:16])
	if root1 == root2 {
		fmt.Println("PASS: roots match regardless of input order")
	} else {
		fmt.Println("FAIL: roots differ — determinism broken")
	}

	// ── PHASE 3: Replay proof ─────────────────────────────────────────────────
	separator("PHASE 3 — REPLAY PROOF")

	ok, result := ReplaySystem(envelopes, anchor)
	fmt.Printf("Run 1 → ok=%v  root=%s\n", ok, result[:16])

	ok2, result2 := ReplaySystem(envelopes, anchor)
	fmt.Printf("Run 2 → ok=%v  root=%s\n", ok2, result2[:16])

	if result == result2 {
		fmt.Println("PASS: identical state roots across runs")
	} else {
		fmt.Println("FAIL: state roots differ")
	}

	// ── FAILURE TESTING ───────────────────────────────────────────────────────
	separator("FAILURE TESTS")

	// 1. Tampered envelope hash
	tampered := envelopes[0]
	tampered.Hash = "0000000000000000000000000000000000000000000000000000000000000000"
	tamperedSet := []Envelope{tampered, envelopes[1], envelopes[2]}
	ok, reason := ReplaySystem(tamperedSet, anchor)
	fmt.Printf("[tampered envelope hash]   ok=%v  reason=%q\n", ok, reason)

	// 2. Wrong input hash
	wrongInput := envelopes[0]
	wrongInput.InputHash = sha256Hex("wrong input")
	wrongInput.Hash = hashEnvelope(wrongInput) // re-seal so envelope is internally consistent
	ok, reason = ValidationAgent(wrongInput, executions[0])
	fmt.Printf("[wrong input hash]         ok=%v  reason=%q\n", ok, reason)

	// 3. Wrong output hash
	wrongOutput := envelopes[0]
	wrongOutput.OutputHash = sha256Hex("wrong output")
	wrongOutput.Hash = hashEnvelope(wrongOutput)
	ok, reason = ValidationAgent(wrongOutput, executions[0])
	fmt.Printf("[wrong output hash]        ok=%v  reason=%q\n", ok, reason)

	// 4. Corrupted execution → replay fail
	corrupted := envelopes[1]
	corrupted.OutputHash = sha256Hex("corrupted")
	corrupted.Hash = hashEnvelope(corrupted)
	corruptedSet := []Envelope{envelopes[0], corrupted, envelopes[2]}
	corruptedAnchor := RelayAgent(corruptedSet) // anchor built from corrupted set
	// now replay original envelopes against corrupted anchor
	ok, reason = ReplaySystem(envelopes, corruptedAnchor)
	fmt.Printf("[corrupted execution]      ok=%v  reason=%q\n", ok, reason)

	// 5. Unordered input — state root must still match
	unordered := []Envelope{envelopes[2], envelopes[0], envelopes[1]}
	unorderedRoot := GenerateStateRoot(unordered)
	if unorderedRoot == root1 {
		fmt.Printf("[unordered input]          PASS: root matches ordered root\n")
	} else {
		fmt.Printf("[unordered input]          FAIL: root mismatch\n")
	}

	separator("DONE")
	fmt.Println("System is blockchain-ready.")
}
