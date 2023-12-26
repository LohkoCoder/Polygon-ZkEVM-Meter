package gasUsedL12

import (
	"Polygon-ZkEVM-meter/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevm"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strconv"
)

var (
	poeContract polygonzkevm.Polygonzkevm
)

const l2BatchVerificationGasUSed = 2000

func WatchBatches(fromBlock int64, toBlock int64) error {
	ethClientLayer1, err := ethclient.Dial(layer1Network.URL)
	ethClientLayer2, err := ethclient.Dial(layer2Network.URL)
	if err != nil {
		return err
	}
	sequenceBatchesEventSignature := crypto.Keccak256Hash([]byte("SequenceBatches(uint64)"))
	verifyBatchesTrustedAggregatorEventSignature := crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint64,bytes32,address)"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(layer1Network.PolygonPoEAddr)},
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
	}
	logs, err := ethClientLayer1.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}
	sequnceBatchesTitle := []string{"                               txHash                           startBatch endBatch gasUsed cdSize  zSize   cdGasUsed cdComRatio  gasComRatio    gasComRatioV"}
	//0xfc46d2b9bd157bdbac59301c5cbb69924f7c4a406f8a8a447284bbd34e915431	397346	397363	1739489	127744	54207	1393420		2.54352		19.03626		18.65028
	utils.Save2File("./sequenceBatchesTxs.txt", sequnceBatchesTitle)
	var verifyBatchesTxs []string
	for _, vLog := range logs {
		l1Receipt, err := ethClientLayer1.TransactionReceipt(context.Background(), common.HexToHash(vLog.TxHash.String()))
		if err != nil {
			return err
		}
		layer1Tx, _, err := ethClientLayer1.TransactionByHash(context.Background(), common.HexToHash(vLog.TxHash.String()))
		if err != nil {
			return err
		}
		switch vLog.Topics[0].String() {
		case sequenceBatchesEventSignature.String():
			// methodId 4Bytes; coinBase 1slot=32Bytes
			callDataSize := uint64(len(layer1Tx.Data())) - 4 - 32
			abi, _ := polygonzkevm.PolygonzkevmMetaData.GetAbi()
			method, err := abi.MethodById(layer1Tx.Data()[:4])
			if err != nil {
				return err
			}

			data, err := method.Inputs.Unpack(layer1Tx.Data()[4:])
			//coinbase := (data[1]).(common.Address)
			var batches []polygonzkevm.PolygonZkEVMBatchData
			bytedata, err := json.Marshal(data[0])
			err = json.Unmarshal(bytedata, &batches)
			endBatch := hex.DecodeUint64(vLog.Topics[1].String())
			startBatch := endBatch - uint64(len(batches)) + 1
			coinbaseZeroBytesNum := analyzeZeroBytesNum(data[1].(common.Address).Bytes()) + 12
			txDataZeroBytesNum := analyzeZeroBytesNum(layer1Tx.Data()[4:])
			callDataZeroBytesNum := txDataZeroBytesNum - coinbaseZeroBytesNum
			callDataGasUsed := uint64(callDataZeroBytesNum)*4 + (callDataSize-uint64(callDataZeroBytesNum))*16

			seqTxStatic := vLog.TxHash.String() + "\t" + strconv.FormatUint(startBatch, 10) + "\t" + strconv.FormatUint(endBatch, 10) +
				"\t" + strconv.FormatUint(l1Receipt.GasUsed, 10) + "\t" + strconv.FormatUint(callDataSize, 10)
			seqTxStatic += "\t" + strconv.FormatUint(uint64(callDataZeroBytesNum), 10) + "\t" + strconv.FormatUint(callDataGasUsed, 10)

			var (
				l2SizeSum                 float64
				l2GasUsedSum              float64
				l1VerifyBatchesGasUsedSum float64
			)

			for i, batch := range batches {
				//fmt.Println(k)
				l1VerifyBatchesGasUsedSum += l2BatchVerificationGasUSed
				l2Txs, err := utils.ParseBatchFromL1(hex.EncodeToString(batch.Transactions))
				if err != nil {
					return err
				}
				var l2TxsBatch []string
				for _, l2Tx := range l2Txs {
					l2Receipt, err := ethClientLayer2.TransactionReceipt(context.Background(), l2Tx.Hash())
					if err != nil {
						return err
					}
					l2TxsBatch = append(l2TxsBatch, l2Tx.Hash().Hex()+"\t"+strconv.FormatUint(l2Receipt.GasUsed, 10)+"\t"+l2Receipt.Size().String()+"\t"+strconv.FormatUint(startBatch+uint64(i), 10))
					l2SizeSum += float64(l2Receipt.Size())
					l2GasUsedSum += float64(l2Receipt.GasUsed)
				}
				utils.Save2File("./TxsGasUsedAndSizeInL2.txt", l2TxsBatch)
			}
			seqTxStatic += "\t\t" + fmt.Sprintf("%.5f", l2SizeSum/float64(callDataSize))
			seqTxStatic += "\t\t" + fmt.Sprintf("%.5f", l2GasUsedSum/float64(l1Receipt.GasUsed))
			seqTxStatic += "\t\t" + fmt.Sprintf("%.5f", l2GasUsedSum/(float64(l1Receipt.GasUsed)+l1VerifyBatchesGasUsedSum))
			utils.Save2File("./sequenceBatchesTxs.txt", []string{seqTxStatic})
		case verifyBatchesTrustedAggregatorEventSignature.String():
			verifyBatchesTxs = append(verifyBatchesTxs, vLog.TxHash.String()+"\t"+strconv.FormatUint(hex.DecodeUint64(vLog.Topics[1].String()), 10)+"\t"+vLog.Topics[2].String()+"\t"+strconv.FormatUint(l1Receipt.GasUsed, 10))
		}
	}
	if len(verifyBatchesTxs) > 0 {
		err = utils.Save2File("./verifyBatchesTxs.txt", verifyBatchesTxs)
		if err != nil {
			return err
		}
	}
	return nil
}

func monitorBatches() error {
	ethClientIns, err := ethclient.Dial(layer1Network.URL)
	if err != nil {
		return err
	}
	sequenceBatchesEventSignature := crypto.Keccak256Hash([]byte("SequenceBatches(uint64)"))
	verifyBatchesTrustedAggregatorEventSignature := crypto.Keccak256Hash([]byte("VerifyBatchesTrustedAggregator(uint64,bytes32,address)"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(layer1Network.PolygonPoEAddr)},
	}
	logs := make(chan types.Log)
	sub, err := ethClientIns.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return err
	}

	for {
		select {
		case err := <-sub.Err():
			return err
		case vLog := <-logs:
			switch vLog.Topics[0].String() {
			case sequenceBatchesEventSignature.String():
				utils.Save2File("./sequenceBatchesTxs.txt", []string{vLog.TxHash.String() + "\t" + strconv.FormatUint(hex.DecodeUint64(vLog.Topics[1].String()), 10)})
			case verifyBatchesTrustedAggregatorEventSignature.String():
				utils.Save2File("./verifyBatchesTxs.txt", []string{vLog.TxHash.String() + "\t" + strconv.FormatUint(hex.DecodeUint64(vLog.Topics[1].String()), 10) + "\t" + vLog.Topics[2].String()})
			}
		}
	}
}

func Layer1BatchDataAnalysis(fromBlock int64, toBlock int64) error {
	ethClientLayer1, err := ethclient.Dial(layer1Network.URL)
	if err != nil {
		return err
	}
	sequenceBatchesEventSignature := crypto.Keccak256Hash([]byte("SequenceBatches(uint64)"))
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(layer1Network.PolygonPoEAddr)},
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
	}
	logs, err := ethClientLayer1.FilterLogs(context.Background(), query)
	if err != nil {
		return err
	}
	for _, vLog := range logs {
		layer1Tx, _, err := ethClientLayer1.TransactionByHash(context.Background(), common.HexToHash(vLog.TxHash.String()))
		if err != nil {
			return err
		}
		switch vLog.Topics[0].String() {
		case sequenceBatchesEventSignature.String():
			fmt.Printf("tx data length %v\n", len(layer1Tx.Data()))
			var zeroByte int
			for _, elem := range layer1Tx.Data() {
				if elem == byte(0) {
					zeroByte++
					//fmt.Printf("zero byte num%v\t", zeroByte)
				}
			}
			fmt.Printf("zero byte num %v\n", zeroByte)
			fmt.Printf("\nzero byte ratio %v\n", zeroByte/len(layer1Tx.Data()))
		}
	}
	return nil
}

func analyzeZeroBytesNum(txDataOrArg []byte) int {
	var zeroByte int
	for _, elem := range txDataOrArg {
		if elem == byte(0) {
			zeroByte++
		}
	}
	return zeroByte
}
