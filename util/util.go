package util

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

//Decode eth Transaction Data
func DecodeRawTx(rawTx string) (*types.Transaction, error) {
	body, err := hexutil.Decode(rawTx)
	if err != nil {
		return nil, err
	}
	var etx types.Transaction
	err = rlp.DecodeBytes(body, &etx)
	if err != nil {
		return nil, err
	}

	return &etx, nil
}

func Uint64ToHexString(val uint64) string {
	return hexutil.EncodeUint64(val)
}

func HexToUint64(hxs string) (uint64, error) {
	return hexutil.DecodeUint64(hxs)
}

func zeroBytes() []byte {
	return []byte{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
}

//parse eth signature
func parseEthSignature(ethtx *types.Transaction) []byte {
	big8 := big.NewInt(8)
	v, r, s := ethtx.RawSignatureValues()
	fmt.Println("v 1:", v.Uint64())
	v = new(big.Int).Sub(v, new(big.Int).Mul(ethtx.ChainId(), big.NewInt(2)))
	fmt.Println("v 2:", v.Uint64())

	rBytes := r.Bytes()
	fmt.Println("rBytes len:", len(rBytes))
	sBytes := s.Bytes()
	fmt.Println("sBytes len:", len(sBytes))
	fmt.Println("tx type:", ethtx.Type())

	var vBytes byte
	if ethtx.Type() == types.LegacyTxType {
		if n := len(rBytes); n < 32 {
			rBytes = append(zeroBytes()[:32-n], rBytes...)
		}

		if n := len(sBytes); n < 32 {
			sBytes = append(zeroBytes()[:32-n], sBytes...)
		}
		v.Sub(v, big8)
		vBytes = byte(v.Uint64() - 27)
	} else if ethtx.Type() == types.DynamicFeeTxType {
		/* V := new(big.Int).Add(v, big.NewInt(27))
		vBytes = V.Bytes() */
	}

	fmt.Println("v 3:", v.Uint64())
	var sign []byte
	sign = append(sign, rBytes...)
	sign = append(sign, sBytes...)
	sign = append(sign, vBytes)
	fmt.Println("sign len:", len(sign))
	return sign
}

//Verify Eth Signature
func VerifyEthSignature(ethtx *types.Transaction) error {
	sign := parseEthSignature(ethtx)
	if len(sign) != 65 {
		return fmt.Errorf("eth signature lenght error:%v", len(sign))
	}

	var sighash common.Hash
	if ethtx.Type() == types.LegacyTxType {
		signer := types.NewEIP155Signer(ethtx.ChainId())
		sighash = signer.Hash(ethtx)
	}
	if ethtx.Type() == types.DynamicFeeTxType {
		signer := types.NewLondonSigner(ethtx.ChainId())
		sighash = signer.Hash(ethtx)
	}

	pub, err := crypto.Ecrecover(sighash[:], sign)
	if err != nil {
		fmt.Println("Ecrecover error:", err)
		return err
	}

	{
		/* msg, _ := ethtx.AsMessage(signer, nil)

		p, err := crypto.SigToPub(ethtx.Hash().Bytes(), sign) //sighash[:], sign)
		if err != nil {
			fmt.Println("SigToPub error:", err)
			return err
		}

		addr := crypto.PubkeyToAddress(*p)
		if msg.From() != addr {
			return fmt.Errorf("verify sender failed! want:%v,got:%v", addr, msg.From())
		}
		fmt.Println("sender:", addr) */
	}

	if !crypto.VerifySignature(pub, sighash[:], sign[:len(sign)-1]) {
		return fmt.Errorf("%v", "Verify Eth Signature failed")
	}
	return nil
}
