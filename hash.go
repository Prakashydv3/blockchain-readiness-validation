package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

func sha256Hex(data string) string {
	h := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", h)
}

func hashEnvelope(e Envelope) string {
	return sha256Hex(e.ExecutionID + e.InputHash + e.OutputHash)
}

func hashAnchor(stateRoot string, timestamp int64) string {
	return sha256Hex(stateRoot + strconv.FormatInt(timestamp, 10))
}
