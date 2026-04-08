package main

// Execution represents a single TTS execution unit
type Execution struct {
	ID     string
	Input  string
	Output string
}

// Envelope is the sealed, hashed record of one execution
type Envelope struct {
	ExecutionID string
	InputHash   string
	OutputHash  string
	Hash        string // SHA256(ExecutionID + InputHash + OutputHash)
}

// Anchor is the L1 commitment record
type Anchor struct {
	StateRoot string
	Timestamp int64
	Hash      string // SHA256(StateRoot + Timestamp)
}
