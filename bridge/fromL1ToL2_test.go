package bridge

import (
	"context"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestBridgeCustomizedUserGasTokenFromL1ToL2(t *testing.T) {
	fromBalance, toBalance, bridgeBalance := checkLayer1AccBalance(true, true)
	t.Logf("L1 before BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalance, toBalance, bridgeBalance)

	fromBalanceL2, toBalanceL2, bridgeBalanceL2 := checkLayer2AccBalance(true)
	t.Logf("L2 before BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalanceL2, toBalanceL2, bridgeBalanceL2)

	BridgeCustomizedUserGasTokenFromL1ToL2()
	time.Sleep(time.Second * 30)
	fromBalance2, toBalance2, bridgeBalance2 := checkLayer1AccBalance(true, true)
	t.Logf("after BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalance2, toBalance2, bridgeBalance2)
	fromBalance2L2, toBalance2L2, bridgeBalance2L2 := checkLayer2AccBalance(true)
	t.Logf("L2 after BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalance2L2, toBalance2L2, bridgeBalance2L2)

	assert.Equal(t, toBalance, toBalance2)
	assert.Equal(t, fromBalanceL2, fromBalance2L2)

	x := big.NewInt(0)
	x.Sub(fromBalance, fromBalance2)
	y := big.NewInt(0)
	y.Sub(toBalance2L2, toBalanceL2)
	z := big.NewInt(0)
	z = z.Sub(bridgeBalance2, bridgeBalance)
	t.Logf("x %v, y%v, z%v\n", x, y, z)
	//assert.Equal(t, x, z)
	//assert.Equal(t, y, z)
}

func TestBridgeEthFromL1ToL2(t *testing.T) {
	fromBalance, toBalance, bridgeBalance := checkLayer1AccBalance(false, true)
	fromBalanceL2, toBalanceL2, bridgeBalanceL2 := checkLayer2AccBalance(true)

	//if fromBalance.Cmp(big.NewInt(int64(math.Pow(10, 17)))) < 1 {
	//	senderPrivateKey, err := crypto.HexToECDSA(strings.TrimPrefix(layer2Network.PrivateKey, "0x"))
	//	chkErr(err)
	//	faucetEth(authLayer2, authLayer1, senderPrivateKey)
	//	nonce, err := layer1Client.PendingNonceAt(context.Background(), authLayer2.From)
	//	authLayer2.Nonce.Set(big.NewInt(int64(nonce + 1)))
	//	fromBalance, _, _ = checkLayer1AccBalance(false, true)
	//}

	BridgeEthFromL1ToL2()
	time.Sleep(time.Second * 30)

	fromBalanceAfter, toBalanceAfter, bridgeBalanceAfter := checkLayer1AccBalance(false, true)
	t.Logf("L1 before BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalance, toBalance, bridgeBalance)
	t.Logf("L1 after  BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalanceAfter, toBalanceAfter, bridgeBalanceAfter)
	fromBalance2L2, toBalance2L2, bridgeBalance2L2 := checkLayer2AccBalance(true)
	t.Logf("L2 before BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalanceL2, toBalanceL2, bridgeBalanceL2)
	t.Logf("L2 after  BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalance2L2, toBalance2L2, bridgeBalance2L2)

	assert.Equal(t, toBalance, toBalanceAfter)
	assert.Equal(t, fromBalanceL2, fromBalance2L2)

	x := big.NewInt(0)
	x.Sub(fromBalance, fromBalanceAfter)
	y := big.NewInt(0)
	y.Sub(toBalance2L2, toBalanceL2)
	z := big.NewInt(0)
	z = z.Sub(bridgeBalanceAfter, bridgeBalance)
	t.Logf("x %v, y%v, z%v\n", x, y, z)
	assert.Equal(t, y, z)
}

func TestBridgeCustomizedUserGasTokenFromL2ToL1(t *testing.T) {
	fromBalance, toBalance, bridgeBalance := checkLayer1AccBalance(true, false)
	t.Logf("L1 before BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalance, toBalance, bridgeBalance)

	fromBalanceL2, toBalanceL2, bridgeBalanceL2 := checkLayer2AccBalance(true)
	t.Logf("L2 before BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalanceL2, toBalanceL2, bridgeBalanceL2)

	BridgeCustomizedUserGasTokenFromL1ToL2()
	time.Sleep(time.Second * 30)
	fromBalance2, toBalance2, bridgeBalance2 := checkLayer1AccBalance(true, false)
	t.Logf("after BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalance2, toBalance2, bridgeBalance2)
	fromBalance2L2, toBalance2L2, bridgeBalance2L2 := checkLayer2AccBalance(true)
	t.Logf("L2 after BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalance2L2, toBalance2L2, bridgeBalance2L2)

	assert.Equal(t, toBalance, toBalance2)
	assert.Equal(t, fromBalanceL2, fromBalance2L2)

	x := big.NewInt(0)
	x.Sub(fromBalance, fromBalance2)
	y := big.NewInt(0)
	y.Sub(toBalance2, toBalance)
	z := big.NewInt(0)
	z = z.Sub(bridgeBalance2, bridgeBalance)
	t.Logf("x %v, y%v, z%v\n", x, y, z)
	//assert.Equal(t, x, z)
}

func TestChkBalance(t *testing.T) {
	t.Logf("acc1 %v, acc2 %v, acc3 %v\n", authLayer1.From.String(), authLayer2.From.String(), authLayer3.From)
	//senderPrivateKey, _ := crypto.HexToECDSA(strings.TrimPrefix(layer2Network.PrivateKey, "0x"))
	//fromBalance2L2, toBalance2L2, bridgeBalance2L2 := checkLayer2AccBalance(true)
	//t.Logf("L2 after BridgeCustomizedUserGasTokenFromL1ToL2 fromBalance %v, toBalance %v, bridgeBalance %v\n", fromBalance2L2, toBalance2L2, bridgeBalance2L2)
	BridgeEthToChargeAuth3()
	bal, _ := layer2Client.PendingBalanceAt(context.Background(), authLayer3.From)
	t.Logf("bal %v\n", bal)
}
