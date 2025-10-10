package keeper_test

import (
    "context"
    "testing"

    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/stretchr/testify/require"

    "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/keeper"
    "github.com/ZK443/qubetics-improvement-pack/chain/x/bridge/types"
)

// Простая заглушка проверяющая, что binding и proof совпадают побайтово.
// В реале — вызывать light/zk верификатор.
type dummyVerifier struct{}
func (d dummyVerifier) VerifyBinding(binding []byte, proof []byte) bool {
    if len(binding) != len(proof) { return false }
    for i := range binding {
        if binding[i] != proof[i] { return false }
    }
    return true
}

func testContext() (sdk.Context, keeper.Keeper) {
    // Для простоты: пустой context. В реале — моки KVStore/keys.
    ctx := sdk.Context{} // минимально для иллюстрации
    k := keeper.NewKeeper(nil, dummyVerifier{})
    return ctx, k
}

func TestVerifyAndExecuteHappyPath(t *testing.T) {
    ctx, k := testContext()

    signer := sdk.AccAddress([]byte("relayer1"))
    executor := sdk.AccAddress([]byte("executor1"))

    msgID := "msg-1"
    digest := []byte("digest")
    binding := keeper.ComputeBindingHash(msgID, digest)

    // Verify
    vresp, err := k.Verify(context.Background(), types.MsgVerifyProof{
        Signer:         signer,
        ProofId:        "proof-1",
        ProofData:      binding,
        MessageID:      msgID,
        ExpectedDigest: digest,
        Verifier:       "dummy",
    })
    require.NoError(t, err)
    require.Equal(t, "verified", vresp.Status)

    // Execute
    eresp, err := k.Execute(context.Background(), types.MsgExecute{
        Executor:  executor,
        MessageId: msgID,
        ProofId:   "proof-1",
        Amount:    sdk.NewInt(1),
    })
    require.NoError(t, err)
    require.Equal(t, "done", eresp.Status)
}

func TestExecuteWithoutVerifyFails(t *testing.T) {
    ctx, k := testContext()
    executor := sdk.AccAddress([]byte("executor1"))
    _, err := k.Execute(context.Background(), types.MsgExecute{
        Executor:  executor,
        MessageId: "msg-no-verify",
        ProofId:   "missing",
        Amount:    sdk.NewInt(1),
    })
    require.ErrorIs(t, err, types.ErrNotVerified)
    _ = ctx
}

func TestReplayExecuteFails(t *testing.T) {
    _, k := testContext()
    signer := sdk.AccAddress([]byte("relayer1"))
    executor := sdk.AccAddress([]byte("executor1"))

    msgID := "msg-2"
    digest := []byte("digest")
    binding := keeper.ComputeBindingHash(msgID, digest)

    _, err := k.Verify(context.Background(), types.MsgVerifyProof{
        Signer:         signer,
        ProofId:        "proof-2",
        ProofData:      binding,
        MessageID:      msgID,
        ExpectedDigest: digest,
        Verifier:       "dummy",
    })
    require.NoError(t, err)

    // Первый Execute — OK
    _, err = k.Execute(context.Background(), types.MsgExecute{
        Executor:  executor,
        MessageId: msgID,
        ProofId:   "proof-2",
        Amount:    sdk.NewInt(1),
    })
    require.NoError(t, err)

    // Повторный Execute — должен упасть
    _, err = k.Execute(context.Background(), types.MsgExecute{
        Executor:  executor,
        MessageId: msgID,
        ProofId:   "proof-2",
        Amount:    sdk.NewInt(1),
    })
    require.ErrorIs(t, err, types.ErrAlreadyExecuted)
}
