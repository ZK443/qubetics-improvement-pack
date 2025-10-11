package types

import "testing"

func TestDefaultParams_Validate(t *testing.T) {
	p := DefaultParams()
	if err := p.Validate(); err != nil {
		t.Fatalf("default params should be valid: %v", err)
	}
}

func TestParams_Validate_Errors(t *testing.T) {
	p := DefaultParams()
	p.RateLimitAmount = 0
	if err := p.Validate(); err == nil {
		t.Fatalf("expected error for zero amount")
	}
	p = DefaultParams()
	p.RateLimitWindowMs = 50
	if err := p.Validate(); err == nil {
		t.Fatalf("expected error for too small window")
	}
}
