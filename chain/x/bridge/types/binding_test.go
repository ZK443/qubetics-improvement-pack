package types

import "testing"

func TestComputeBindingHash_Differs(t *testing.T) {
	a := ComputeBindingHash("id-1", []byte("m1"))
	b := ComputeBindingHash("id-2", []byte("m1"))
	if a == b {
		t.Fatalf("hashes should differ for different msgIDs")
	}
}
