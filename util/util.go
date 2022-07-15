package util

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"toprelayer/config"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/ssh/terminal"

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
	v, r, s := ethtx.RawSignatureValues()

	rBytes := r.Bytes()
	if n := len(rBytes); n < 32 {
		rBytes = append(zeroBytes()[:32-n], rBytes...)
	}

	sBytes := s.Bytes()
	if n := len(sBytes); n < 32 {
		sBytes = append(zeroBytes()[:32-n], sBytes...)
	}

	var vBytes byte
	big8 := big.NewInt(8)
	if ethtx.Type() == types.LegacyTxType {
		v = new(big.Int).Sub(v, new(big.Int).Mul(ethtx.ChainId(), big.NewInt(2)))
		v.Sub(v, big8)
		vBytes = byte(v.Uint64() - 27)
	}

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

	if !crypto.VerifySignature(pub, sighash[:], sign[:len(sign)-1]) {
		return fmt.Errorf("%v", "Verify Eth Signature failed")
	}
	return nil
}

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	pass, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return string(pass), err
	}
	fmt.Println()
	return string(pass), nil
}

func GetPasses(cfg *config.Config) (map[string]string, error) {
	chainpass := make(map[string]string)
	for _, chain := range cfg.RelayerConfig {
		if !chain.Start {
			continue
		}
		pass, err := readPassword(">>> Please Enter " + chain.Chain + " pasword:\n>>> ")
		if err != nil {
			log.Fatal("Read password file failed:", err)
			return nil, err
		}
		chainpass[chain.Chain] = pass
	}
	return chainpass, nil
}
