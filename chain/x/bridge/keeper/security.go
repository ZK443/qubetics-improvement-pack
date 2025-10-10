package keeper

import (
    "crypto/sha256"
)

// ComputeBindingHash(msgID || expectedDigest)
func ComputeBindingHash(messageID string, expectedDigest []byte) []byte {
    h := sha256.Sum256(append([]byte(messageID), expectedDigest...))
    return h[:]
}
