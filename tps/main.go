package tps

import (
	"Polygon-ZkEVM-meter/utils"
	"context"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"math/big"
	"strings"
	"time"
)

const (
	txTimeout = 60 * time.Second
)

func genAcc() {
	receiverKey, err := crypto.GenerateKey()
	chkErr(err)
	receiverAddr := crypto.PubkeyToAddress(receiverKey.PublicKey)
	log.Infof("receiverPriKey %v, receiverAddr %v\n", hex.EncodeToString(crypto.FromECDSA(receiverKey)), receiverAddr)

}

func OneHundredTxs() {
	network := utils.Networks[0]
	ctx := context.Background()

	log.Infof("connecting to %v: %v", network.Name, network.URL)
	client, err := ethclient.Dial(network.URL)
	chkErr(err)
	log.Infof("connected")

	receiverKey, err := crypto.GenerateKey()
	chkErr(err)
	log.Infof("receiverPriKey %v\n", hex.EncodeToString(crypto.FromECDSA(receiverKey)))
	receiverAddr := crypto.PubkeyToAddress(receiverKey.PublicKey)
	receiverNonce, err := client.PendingNonceAt(ctx, receiverAddr)

	//var receiverAddr = common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D")
	auth := operations.MustGetAuth(network.PrivateKey, network.ChainID)
	chkErr(err)
	balance, err := client.BalanceAt(ctx, auth.From, nil)
	chkErr(err)
	//24999998517100000017320035
	//24999998514900000025986225
	log.Debugf("ETH Balance for %v: %v", auth.From, balance)
	senderNonce, err := client.PendingNonceAt(ctx, auth.From)
	chkErr(err)

	senderPrivateKey, err := crypto.HexToECDSA(strings.TrimPrefix(network.PrivateKey, "0x"))
	chainid := big.NewInt(int64(network.ChainID))

	faucetTx := types.NewTx(&types.LegacyTx{
		Nonce:    senderNonce,
		To:       &receiverAddr,
		Value:    big.NewInt(int64(math.Pow(10, 18))),
		GasPrice: big.NewInt(int64(math.Pow(10, 9))),
		Gas:      31000,
	})
	signedTx, err := types.SignTx(faucetTx, types.NewEIP155Signer(chainid), senderPrivateKey)
	err = client.SendTransaction(ctx, signedTx)
	chkErr(err)
	time.Sleep(time.Second * 5)
	//err = operations.WaitTxToBeMined(ctx, client, faucetTx, 1*operations.DefaultTimeoutTxToBeMined)
	//chkErr(err)
	receiverBalance, err := client.BalanceAt(ctx, receiverAddr, nil)
	chkErr(err)
	log.Debugf("ETH Balance for %v: %v", receiverAddr, receiverBalance)

	senderNonce++

	var txs []*types.Transaction
	for i := 0; i < 50; i++ {
		basicTx := types.NewTx(&types.LegacyTx{
			Nonce:    senderNonce,
			To:       &receiverAddr,
			Value:    big.NewInt(12345),
			GasPrice: big.NewInt(int64(math.Pow(10, 9))),
			Gas:      31000,
		})
		signedTx, err := types.SignTx(basicTx, types.NewEIP155Signer(chainid), senderPrivateKey)
		chkErr(err)
		txs = append(txs, signedTx)
		senderNonce++

		basicTx2 := types.NewTx(&types.LegacyTx{
			Nonce:    receiverNonce,
			To:       &auth.From,
			Value:    big.NewInt(123),
			GasPrice: big.NewInt(int64(math.Pow(10, 9))),
			Gas:      31000,
		})
		signedTx, err = types.SignTx(basicTx2, types.NewEIP155Signer(chainid), receiverKey)
		chkErr(err)
		txs = append(txs, signedTx)
		receiverNonce++

		//basicTx3 := types.NewTx(&types.LegacyTx{
		//	Nonce:    senderNonce,
		//	To:       &receiverAddr,
		//	Value:    big.NewInt(12),
		//	GasPrice: big.NewInt(int64(math.Pow(10, 9))),
		//	Gas:      31000,
		//})
		//signedTx, err = types.SignTx(basicTx3, types.NewEIP155Signer(chainid), receiverKey)
		//chkErr(err)
		//txs = append(txs, signedTx)
		//senderNonce++
	}

	log.Infof("1", time.Now())
	for i := 0; i < 100; i++ {
		log.Infof("tx NO. %v\n", i)
		err = client.SendTransaction(ctx, txs[i])
		chkErr(err)
	}

	err = operations.WaitTxToBeMined(ctx, client, txs[99], 1*operations.DefaultTimeoutTxToBeMined)
	chkErr(err)
	log.Infof("2", time.Now())
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
