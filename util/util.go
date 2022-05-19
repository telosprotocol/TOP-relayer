package util

import (
	"fmt"
	"math/big"

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
	v = new(big.Int).Sub(v, new(big.Int).Mul(ethtx.ChainId(), big.NewInt(2)))
	v.Sub(v, big8)

	rBytes := r.Bytes()
	if n := len(rBytes); n < 32 {
		rBytes = append(zeroBytes()[:32-n], rBytes...)
	}
	sBytes := s.Bytes()
	if n := len(sBytes); n < 32 {
		sBytes = append(zeroBytes()[:32-n], sBytes...)
	}
	vBytes := byte(v.Uint64() - 27)

	var sign []byte
	sign = append(sign, rBytes...)
	sign = append(sign, sBytes...)
	sign = append(sign, vBytes)
	return sign
}

//Verify Eth Signature
func VerifyEthSignature(ethtx *types.Transaction) error {
	sign := parseEthSignature(ethtx)
	if len(sign) != 65 {
		return fmt.Errorf("eth signature lenght error:%v", len(sign))
	}

	signer := types.NewEIP155Signer(ethtx.ChainId())
	sighash := signer.Hash(ethtx)
	pub, err := crypto.Ecrecover(sighash[:], sign)
	if err != nil {
		return err
	}

	{

		signer := types.NewEIP155Signer(ethtx.ChainId())
		msg, _ := ethtx.AsMessage(signer, nil)

		p, err := crypto.SigToPub(sighash[:], sign)
		if err != nil {
			return err
		}

		addr := crypto.PubkeyToAddress(*p)
		if msg.From() != addr {
			return fmt.Errorf("verify sender failed! want:%v,got:%v", addr, msg.From())
		}
	}

	if !crypto.VerifySignature(pub, sighash[:], sign[:64]) {
		return fmt.Errorf("%v", "Verify Eth Signature failed")
	}
	return nil
}
