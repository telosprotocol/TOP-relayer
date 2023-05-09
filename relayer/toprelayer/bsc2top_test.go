package toprelayer

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
	"testing"
	"toprelayer/config"
	"toprelayer/relayer/toprelayer/congress"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/wonderivan/logger"
)

const bscUrl = "https://bsc-dataseed4.binance.org"

func TestGetBscInitData(t *testing.T) {
	var height uint64 = 20250000

	ethclient, err := ethclient.Dial(hecoUrl)
	if err != nil {
		t.Fatal(err)
	}
	var batch []byte
	for i := height; i <= height+11; i++ {
		header, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(i))
		if err != nil {
			t.Fatal(err)
		}
		rlp_bytes, err := rlp.EncodeToBytes(header)
		if err != nil {
			t.Fatal(err)
		}
		batch = append(batch, rlp_bytes...)
	}
	logger.Debug(common.Bytes2Hex(batch))
}

func TestGetBscSyncData(t *testing.T) {
	var start_height uint64 = 17276022
	var sync_num uint64 = 1

	ethsdk, err := ethclient.Dial(hecoUrl)
	if err != nil {
		t.Fatal(err)
	}
	con := congress.New(ethsdk)
	err = con.Init(start_height - 1)
	if err != nil {
		t.Fatal(err)
	}

	var batch []byte
	for h := start_height; h <= start_height+sync_num-1; h++ {
		header, err := ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(17276012))
		if err != nil {
			t.Fatal(err)
		}
		out, err := con.GetLastSnapBytes(header)
		if err != nil {
			t.Fatal(err)
		}
		batch = append(batch, out...)
	}
	logger.Debug(common.Bytes2Hex(batch))
}

func TestBscInit(t *testing.T) {
	// changable
	var start_height uint64 = 20250000
	var sync_num uint64 = 12
	var topUrl string = "http://192.168.30.200:8080"
	var keyPath = "../../.relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     []string{topUrl},
		KeyPath: keyPath,
	}
	relayer := &Heco2TopRelayer{}
	err := relayer.Init(cfg, []string{bscUrl}, defaultPass)
	if err != nil {
		t.Fatal(err)
	}
	var batch []byte
	for h := start_height; h <= start_height+sync_num-1; h++ {
		header, err := relayer.ethsdk.HeaderByNumber(context.Background(), big.NewInt(0).SetUint64(h))
		if err != nil {
			t.Fatal(err)
		}
		rlp_bytes, err := rlp.EncodeToBytes(header)
		if err != nil {
			t.Fatal("EncodeToBytes: ", err)
		}
		batch = append(batch, rlp_bytes...)
	}

	nonce, err := relayer.wallet.NonceAt(context.Background(), relayer.wallet.Address(), nil)
	if err != nil {
		t.Fatal(err)
	}
	gaspric, err := relayer.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	ops := &bind.TransactOpts{
		From:      relayer.wallet.Address(),
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  500000,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    relayer.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}
	tx, err := relayer.transactor.Init(ops, batch, string(""))
	if err != nil {
		t.Fatal(err)
	}
	logger.Debug(tx.Hash())
}

func TestHecoDemo(t *testing.T) {
	// top rpc地址
	var topUrl string = "http://192.168.95.3:8080"
	// keystore 和账户是对应的，这里我写的我的绝对路径 /Users/pzxy/Workspace/Top/TOP-relayer/.relayer/wallet/top
	var keyPath = ".relayer/wallet/top"

	cfg := &config.Relayer{
		Url:     []string{topUrl},
		KeyPath: keyPath,
	}
	relayer := &Heco2TopRelayer{}
	err := relayer.Init(cfg, []string{hecoUrl}, defaultPass)
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := relayer.wallet.NonceAt(context.Background(), relayer.wallet.Address(), nil)
	if err != nil {
		t.Error(err)
	}
	gaspric, err := relayer.wallet.SuggestGasPrice(context.Background())
	if err != nil {
		t.Error(err)
	}
	//must init ops as bellow
	ops := &bind.TransactOpts{
		From:      relayer.wallet.Address(),
		Nonce:     big.NewInt(0).SetUint64(nonce),
		GasLimit:  50000,
		GasFeeCap: gaspric,
		GasTipCap: big.NewInt(0),
		Signer:    relayer.signTransaction,
		Context:   context.Background(),
		NoSend:    false,
	}
	// 将 1.初始化 获取到的数据填写到这里
	hexData := `将 1.初始化 获取到的数据填写到这里`
	initHeaders := hexutils.HexToBytes(hexData)
	tx, err := relayer.transactor.Init(ops, initHeaders, "")
	if err != nil {
		t.Error(err)
	}
	t.Log(tx.Hash())
}

func TestETH(t *testing.T) {
	// 连接到以太坊网络
	client, err := ethclient.Dial("http://192.168.95.3:8080")
	if err != nil {
		panic(err)
	}

	// 构建交易数据
	privateKey, err := crypto.HexToECDSA("ef2a060a2b4d1661d39c36c876148fa1d8952fb3250507df231ebe7bc95b16a0")
	if err != nil {
		panic(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress("0x912eE24f30fd54E5b83Bc181Bb56a767B68ba086"))
	if err != nil {
		panic(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	toAddress := common.HexToAddress("RECIPIENT_ADDRESS")
	value := big.NewInt(1000000000000000000) // 1 ETH
	gasLimit := uint64(21000)
	data := []byte("")

	// 创建交易对象
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	// 对交易进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(1)), privateKey)
	if err != nil {
		panic(err)
	}

	// 在以太坊网络上广播交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("TxHash: %s\n", signedTx.Hash().Hex())
}

//{"id":3601684829,"jsonrpc":"2.0","method":"eth_estimateGas","params":[{"from":"0x8cb56de7306ece6d9cb0fa4c9ddb623b52b8d509","to":"0xb5964709bb7a9f28369d45172ddb3362b27e9cf3","data":"0xe488d078","value":"0x0"}]}
func TestTOPGas(t *testing.T) {
	// 连接以太坊客户端
	client, err := ethclient.Dial("http://192.168.95.3:8080")
	if err != nil {
		panic(err)
	}

	// 获取当前以太坊网络的GasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress("0xB5964709BB7a9F28369d45172DdB3362b27E9cf3"))
	if err != nil {
		panic(err)
	}
	// 构建合约方法参数
	// 假设合约方法签名为 "demo2()"，不需要传递任何参数
	input, err := hex.DecodeString("e488d078")
	if err != nil {
		panic(err)
	}

	// 创建一个未签名的交易对象
	tx := types.NewTransaction(
		nonce, // nonce值
		common.HexToAddress("0xb5964709bb7a9f28369d45172ddb3362b27e9cf3"), // 合约地址
		big.NewInt(0), // 转账金额为0
		0,             // gasLimit，此处为0，会被 eth_estimateGas 覆盖
		gasPrice,      // gasPrice
		input,         // 合约方法ABI编码
	)

	// 调用eth_estimateGas方法估算燃气成本
	privateKey, err := crypto.HexToECDSA("ef2a060a2b4d1661d39c36c876148fa1d8952fb3250507df231ebe7bc95b16a0")

	callMsg := ethereum.CallMsg{
		From:  crypto.PubkeyToAddress(privateKey.PublicKey),
		To:    tx.To(),
		Data:  tx.Data(),
		Value: big.NewInt(0),
	}
	gasLimit, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 输出估算的燃气成本
	fmt.Printf("Estimated gas limit: %d\n", gasLimit)
}

func TestETHGas(t *testing.T) {
	// 连接以太坊客户端
	client, err := ethclient.Dial("http://192.168.95.3:8080")
	if err != nil {
		panic(err)
	}

	// 获取当前以太坊网络的GasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress("0xB5964709BB7a9F28369d45172DdB3362b27E9cf3"))
	if err != nil {
		panic(err)
	}
	// 构建合约方法参数
	// 假设合约方法签名为 "demo2()"，不需要传递任何参数
	input, err := hex.DecodeString("9e457ea7")
	if err != nil {
		panic(err)
	}

	// 创建一个未签名的交易对象
	tx := types.NewTransaction(
		nonce, // nonce值
		common.HexToAddress("0xB5964709BB7a9F28369d45172DdB3362b27E9cf3"), // 合约地址
		big.NewInt(0), // 转账金额为0
		0,             // gasLimit，此处为0，会被 eth_estimateGas 覆盖
		gasPrice,      // gasPrice
		input,         // 合约方法ABI编码
	)

	// 调用eth_estimateGas方法估算燃气成本
	callMsg := ethereum.CallMsg{
		To:   tx.To(),
		Data: tx.Data(),
	}
	gasLimit, err := client.EstimateGas(context.Background(), callMsg)
	if err != nil {
		fmt.Println(err.Error())
	}

	// 输出估算的燃气成本
	fmt.Printf("Estimated gas limit: %d\n", gasLimit)
}
