package bridge

import (
	"Polygon-ZkEVM-meter/utils"
	"context"
	"crypto/ecdsa"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/matic"
	"github.com/0xPolygonHermez/zkevm-node/etherman/smartcontracts/polygonzkevmbridge"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/test/operations"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math"
	"math/big"
	"strings"
	"time"
)

var (
	layer2Network   = utils.Networks[0]
	layer1Network   = utils.Networks[1]
	layer2Client, _ = ethclient.Dial(layer2Network.URL)
	layer1Client, _ = ethclient.Dial(layer1Network.URL)
	bridgeAddr      = "0x79A227ca1609C40856B52b09b24DDb98ba8F3b7b"
	layer2Bridge, _ = polygonzkevmbridge.NewPolygonzkevmbridge(common.HexToAddress(bridgeAddr), layer2Client)
	layer1Bridge, _ = polygonzkevmbridge.NewPolygonzkevmbridge(common.HexToAddress(bridgeAddr), layer1Client)
	authLayer2      = operations.MustGetAuth(layer2Network.PrivateKey, layer2Network.ChainID)
	authLayer1      = operations.MustGetAuth(layer1Network.PrivateKey, layer1Network.ChainID)
	authLayer3      = operations.MustGetAuth(utils.Networks[2].PrivateKey, layer1Network.ChainID)
	callLayer2Opts  = &bind.CallOpts{
		Pending:     false,
		From:        authLayer2.From,
		BlockNumber: nil,
		Context:     context.Background(),
	}
	callLayer1Opts = &bind.CallOpts{
		Pending:     false,
		From:        authLayer1.From,
		BlockNumber: nil,
		Context:     context.Background(),
	}
	maticAddr     = "0xdB5Db3d21b3Fe51900A8332B8EEE32d675C438f1"
	maticToken, _ = matic.NewMatic(common.HexToAddress(maticAddr), layer1Client)
	txTimeout     = 60 * time.Second
)

func bridgeL1Count() {
	curCount, err := layer1Bridge.DepositCount(callLayer2Opts)
	chkErr(err)
	log.Debugf("curCount %v\n", curCount)
}

func checkL1L2NetworkId() {
	l1BriNetworkId, _ := layer1Bridge.NetworkID(callLayer2Opts)
	log.Debugf("l1 bridge networkId %v\n", l1BriNetworkId)
	l2BriNetworkId, _ := layer2Bridge.NetworkID(callLayer2Opts)
	log.Debugf("l2 bridge networkId %v\n", l2BriNetworkId)
}

func faucetEth(fromAuth *bind.TransactOpts, toAuth *bind.TransactOpts, fromPriKey *ecdsa.PrivateKey) {
	senderNonce, err := layer1Client.PendingNonceAt(context.Background(), fromAuth.From)
	chkErr(err)
	faucetTx := types.NewTx(&types.LegacyTx{
		Nonce:    senderNonce,
		To:       &toAuth.From,
		Value:    big.NewInt(int64(math.Pow(10, 18)) * 100),
		GasPrice: big.NewInt(int64(math.Pow(10, 9))),
		Gas:      31000,
	})
	signedTx, err := types.SignTx(faucetTx, types.NewEIP155Signer(big.NewInt(int64(layer1Network.ChainID))), fromPriKey)
	chkErr(err)
	err = layer1Client.SendTransaction(context.Background(), signedTx)
	time.Sleep(time.Second * 5)
	err = operations.WaitTxToBeMined(context.Background(), layer1Client, signedTx, txTimeout)
	chkErr(err)
}

func faucetCustomizedUserGasToken(fromAuth *bind.TransactOpts, toAuth *bind.TransactOpts) {
	tx, err := maticToken.Transfer(fromAuth, toAuth.From, big.NewInt(int64(math.Pow(10, 18)*15)))
	chkErr(err)
	operations.WaitTxToBeMined(context.Background(), layer1Client, tx, txTimeout)
}

func BridgeEthFromL1ToL2() {
	authLayer1.Value = big.NewInt(int64(math.Pow(10, 5)))
	BridgeEth(authLayer1, 1, authLayer2.From, true, authLayer1.Value)
}

func BridgeEthToChargeAuth3() {
	authLayer2.Value = big.NewInt(int64(math.Pow(10, 18)) * 100)
	BridgeEth(authLayer2, 1, authLayer3.From, true, authLayer2.Value)
}

func BridgeEth(auth *bind.TransactOpts, destinationNetwork uint32, destinationAddr common.Address, isLayer1ToLayer2 bool, amount *big.Int) {
	var bridge *polygonzkevmbridge.Polygonzkevmbridge
	if isLayer1ToLayer2 {
		bridge = layer1Bridge
	} else {
		bridge = layer2Bridge
	}
	tx, err := bridge.BridgeAsset(auth, destinationNetwork, destinationAddr, amount, state.ZeroAddress, true, []byte{})
	chkErr(err)
	err = operations.WaitTxToBeMined(context.Background(), layer1Client, tx, txTimeout)
	chkErr(err)

}

func BridgeCustomizedUserGasToken(auth *bind.TransactOpts, destinationNetwork uint32, destinationAddr common.Address, isLayer1ToLayer2 bool, amount *big.Int) {
	var bridge *polygonzkevmbridge.Polygonzkevmbridge
	if isLayer1ToLayer2 {
		bridge = layer1Bridge
	} else {
		bridge = layer2Bridge
	}

	tx, err := bridge.BridgeAsset(auth, destinationNetwork, destinationAddr, amount, common.HexToAddress(maticAddr), true, []byte{})
	chkErr(err)
	err = operations.WaitTxToBeMined(context.Background(), layer1Client, tx, txTimeout)
	chkErr(err)
}

func BridgeCustomizedUserGasTokenFromL1ToL2() {
	l1SenderMaticBalance, _ := maticToken.BalanceOf(callLayer1Opts, authLayer1.From)

	layer1Balance, _ := layer1Client.PendingBalanceAt(context.Background(), authLayer1.From)
	if layer1Balance.Cmp(big.NewInt(int64(math.Pow(10, 17)))) < 1 {
		senderPrivateKey, err := crypto.HexToECDSA(strings.TrimPrefix(layer2Network.PrivateKey, "0x"))
		chkErr(err)
		faucetEth(authLayer2, authLayer1, senderPrivateKey)
		layer1Balance, _ = layer1Client.PendingBalanceAt(context.Background(), authLayer1.From)
	}

	if l1SenderMaticBalance.Cmp(big.NewInt(int64(math.Pow(10, 18)))) < 1 {
		faucetCustomizedUserGasToken(authLayer2, authLayer1)
	}
	tx, err := maticToken.Approve(authLayer1, common.HexToAddress(bridgeAddr), big.NewInt(int64(math.Pow(10, 18))))
	chkErr(err)
	err = operations.WaitTxToBeMined(context.Background(), layer1Client, tx, txTimeout)

	authLayer1.Value = big.NewInt(0)
	//permitData, err := hex.DecodeString("d505accf00000000000000000000000004e82b9508c751be8c745bfcefe8187c9024eb5c00000000000000000000000079a227ca1609c40856b52b09b24ddb98ba8f3b7b00000000000000000000000000000000000000000000000000038d7ea4c68000000000000000000000000000000000000000000000000000000000006608ab4a000000000000000000000000000000000000000000000000000000000000001cb41024493d4e09f0048a8e00ed68ac8e088d6dbf8fce5d4b3a70849148be02ff6bdd67b72e0250a838de83567ea80130f597dc8f33f61597a858da0173824e58")
	//chkErr(err)
	amount := big.NewInt(int64(math.Pow(10, 1)))
	BridgeCustomizedUserGasToken(authLayer1, 1, authLayer2.From, true, amount)
}

func checkLayer1AccBalance(isCustomziedUserGasToken bool, isLayer1ToLayer2 bool) (fromBalance *big.Int, toBalance *big.Int, bridgeBalance *big.Int) {
	if isCustomziedUserGasToken && isLayer1ToLayer2 {
		fromBalance, _ = maticToken.BalanceOf(callLayer1Opts, authLayer1.From)
		bridgeBalance, _ = maticToken.BalanceOf(callLayer1Opts, common.HexToAddress(bridgeAddr))
		toBalance, _ = maticToken.BalanceOf(callLayer1Opts, authLayer2.From)
		return
	}

	if !isCustomziedUserGasToken && isLayer1ToLayer2 {
		fromBalance, _ = layer1Client.PendingBalanceAt(context.Background(), authLayer1.From)
		toBalance, _ = layer1Client.PendingBalanceAt(context.Background(), authLayer2.From)
		bridgeBalance, _ = layer1Client.PendingBalanceAt(context.Background(), common.HexToAddress(bridgeAddr))
		return
	}

	return nil, nil, nil
}

func checkLayer2AccBalance(isLayer1ToLayer2 bool) (fromBalance *big.Int, toBalance *big.Int, bridgeBalance *big.Int) {

	if isLayer1ToLayer2 {
		fromBalance, _ = layer2Client.PendingBalanceAt(context.Background(), authLayer1.From)
		toBalance, _ = layer2Client.PendingBalanceAt(context.Background(), authLayer2.From)
	} else {
		fromBalance, _ = layer2Client.PendingBalanceAt(context.Background(), authLayer2.From)
		toBalance, _ = layer2Client.PendingBalanceAt(context.Background(), authLayer1.From)
	}
	bridgeBalance, _ = layer2Client.PendingBalanceAt(context.Background(), common.HexToAddress(bridgeAddr))
	return
}

func chkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
