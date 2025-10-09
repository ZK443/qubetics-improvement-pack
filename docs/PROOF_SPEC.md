# Proof Formats (Draft)

## Ethereum (ProofETH)

- `Header`: RLP-encoded block header (or SSZ if using beacon light client).  
- `Receipt`: Merkle proof to the receipt containing the bridge event.  
- `LogIndex`: index of the event in the receipt’s logs.  

**Validation:** verify header → receipt root → log inclusion; check event params ⇢ message hash/nonce/route.

## Solana (ProofSOL)

- `BlockHeader`: last finalized header (light-client path).  
- `MerklePath`: proof of message account/log entry.  
- `Index`: entry index in the tree.  

**Validation:** verify header, then inclusion proof, then payload binding.

## Bitcoin (ProofBTC, SPV)

- `Headers[]`: contiguous header chain with sufficient work.  
- `TxID`: transaction id we expect to see.  
- `MerklePath[]`: branch proving TxID inclusion in the block.  

**Validation:** cumulative work ≥ threshold; merkle path to TxID; script/payload binding.
