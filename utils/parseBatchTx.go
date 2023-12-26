package utils

import (
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/core/types"
)

func ParseBatchFromL1(batchData string) ([]types.Transaction, error) {
	pre155, err := hex.DecodeString(batchData)
	if err != nil {
		return nil, err
	}
	txs, _, _, err := state.DecodeTxs(pre155, 4)
	return txs, err
}
