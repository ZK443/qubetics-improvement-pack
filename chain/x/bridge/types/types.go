package types

// Generic message format and proof container (skeleton).
// In a real impl, proof types would be concrete (light/zk clients).

type Route string
type ChainID string

type Message struct {
	ID        string   // unique id (hash)
	Nonce     uint64   // replay protection
	Source    ChainID  // from chain
	Dest      ChainID  // to chain
	Route     Route    // token-transfer / contract-call / etc
	Payload   []byte   // encoded payload
}

type Proof struct {
	Client    string   // client type (light/zk)
	Data      []byte   // proof bytes
	Header    []byte   // header/commitment for verification
}

type VerificationResult struct {
	Valid   bool
	Reason  string
}

const (
	ModuleName = "bridge"
)
