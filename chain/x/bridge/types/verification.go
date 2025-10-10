package types

// LightVerifier и ZkVerifier — интерфейсы для настоящих клиентов.
// В PoC мы просто делаем заглушку VerifyBinding(binding, proof).
type BindingVerifier interface {
    VerifyBinding(binding []byte, proof []byte) bool
}
