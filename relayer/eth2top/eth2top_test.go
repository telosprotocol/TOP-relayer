package eth2top

import (
	"math/big"
	"testing"
	"toprelayer/base"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const SUBMITTERURL string = "http://192.168.50.204:19086"

const LISTENURL string = "http://192.168.50.235:8545"

var DEFAULTPATH = "../../.relayer/wallet/top"
var CONTRACT common.Address = common.HexToAddress("0xa3e165d80c949833C5c82550D697Ab31Fd3BB446")

func TestSubmitHeader(t *testing.T) {
	sub := &Eth2TopRelayer{}
	err := sub.Init(SUBMITTERURL, LISTENURL, DEFAULTPATH, "", base.TOP, CONTRACT, 90, 0, false)
	if err != nil {
		t.Fatal(err)
	}
	/* 	wg := &sync.WaitGroup{}
	   	err = sub.StartRelayer(wg)
	   	if err != nil {
	   		t.Fatal(err)
	   	} */

	/* 	hashes, err := sub.signAndSendTransactions(1, 10)
	   	if err != nil {
	   		t.Fatal("signAndSendTransactions error:", err)
	   	}
	   	t.Log("hashes:", hashes) */

	nonce, err := sub.wallet.GetNonce(sub.wallet.CurrentAccount().Address)
	if err != nil {
		t.Fatal("GetNonce error:", err)
	}
	balance, err := sub.wallet.GetBalance(sub.wallet.CurrentAccount().Address)
	if err != nil {
		t.Fatal("GetBalance error:", err)
	}
	t.Log("balance:", balance, "nonce:", nonce)

	var headers []*types.Header
	for i := 1; i <= 2; i++ {
		headers = append(headers, &types.Header{Number: big.NewInt(int64(i))})
	}
	hash, err := sub.batch(headers, nonce)
	if err != nil {
		t.Fatal("batch error:", err)
	}
	t.Log("stx hash:", hash)

	/* data, err := msg.EncodeHeaders(&headers)
	if err != nil {
		t.Fatal("EncodeToBytes:", err)
	}

	if sub.wallet == nil {
		t.Fatal("nil wallet!!!")
	}

	stx, err := sub.submitEthHeader(data, nonce)
	if err != nil {
		t.Fatal("SubmitHeader error:", err)
	}
	t.Log("stx hash:", stx.Hash(), "type:", stx.Type())

	byt, err := stx.MarshalBinary()
	if err != nil {
		t.Fatal("MarshalBinary error:", err)
	}
	t.Log("rawtx:", hexutil.Encode(byt)) */
}
