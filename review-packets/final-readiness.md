# Final Readiness Report

## Full Execution Flow

```
Gurukul TTS Execution
  → ExecutionAgent     — seals input/output into Envelope (SHA256 hashes)
  → ValidationAgent    — replays envelope, rejects on any hash mismatch
  → RelayAgent         — collects validated envelopes, generates StateRoot, anchors to L1
  → ReplaySystem       — full deterministic replay verification
```

## State Root Proof

State root is computed by:
1. Sorting envelopes by `ExecutionID` (lexicographic, mandatory)
2. Concatenating all envelope hashes in sorted order
3. SHA256 of the concatenated string

Order-independence proof (live run):

```
A,B,C → 2b9af26af6b80367...
C,B,A → 2b9af26af6b80367...
PASS: roots match regardless of input order
```

Same input always produces the same root. No randomness. No timestamps in state root.

## Agent Behavior

| Agent | Role | Rejects on |
|---|---|---|
| ExecutionAgent | Seals execution into Envelope | — |
| ValidationAgent | Replays envelope against source execution | input hash mismatch, output hash mismatch, envelope hash mismatch |
| RelayAgent | Generates StateRoot + L1 Anchor | — (only receives validated envelopes) |

Live agent log (3 executions):

```
[ExecutionAgent] sealed  id=exec-A  hash=5f294ed628ac
[ValidationAgent] valid  id=exec-A
[ExecutionAgent] sealed  id=exec-B  hash=0a07607b695c
[ValidationAgent] valid  id=exec-B
[ExecutionAgent] sealed  id=exec-C  hash=b7b156bb038d
[ValidationAgent] valid  id=exec-C
[RelayAgent]  state_root=2b9af26af6b8
[RelayAgent]  anchor_hash=5f0f7eb04d73
```

## Replay Proof

`ReplaySystem` verifies:
1. Every envelope hash recomputed from fields matches stored hash
2. StateRoot recomputed from envelopes matches anchor's StateRoot
3. Anchor hash recomputed from StateRoot + Timestamp matches stored anchor hash

Live replay proof:

```
Run 1 → ok=true  root=2b9af26af6b80367
Run 2 → ok=true  root=2b9af26af6b80367
PASS: identical state roots across runs
```

Corrupted execution → replay fails before anchor acceptance:

```
[corrupted execution]  ok=false  reason="state root mismatch: got 2b9af26af6b8 want 72ecf48997f2"
```

## Failure Cases

All failures are caught **before** anchor acceptance:

```
[tampered envelope hash]   ok=false  reason="envelope hash mismatch: exec-A"
[wrong input hash]         ok=false  reason="input hash mismatch for exec-A"
[wrong output hash]        ok=false  reason="output hash mismatch for exec-A"
[corrupted execution]      ok=false  reason="state root mismatch: got 2b9af26af6b8 want 72ecf48997f2"
[unordered input]          PASS: root matches ordered root
```

## Verdict

→ READY
