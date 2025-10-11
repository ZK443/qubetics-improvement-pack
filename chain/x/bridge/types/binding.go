package types

import "crypto/sha256"

// ComputeBindingHash строит детерминированный хеш для привязки доказательства к сообщению.
// Формат: sha256(RouterKey || 0x00 || msgID || 0x00 || meta)
func ComputeBindingHash(msgID string, meta []byte) [32]byte {
	data := make([]byte, 0, len(RouterKey)+1+len(msgID)+1+len(meta))
	data = append(data, []byte(RouterKey)...)
	data = append(data, 0x00)
	data = append(data, []byte(msgID)...)
	data = append(data, 0x00)
	data = append(data, meta...)
	return sha256.Sum256(data)
}
