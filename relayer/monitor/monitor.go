package monitor

import (
	"container/list"
	"context"
	"math/big"
	"time"
	"toprelayer/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/wonderivan/logger"
)

const (
	maxErrorNum = 60

	checkTxInterval      = 5
	checkAccountInterval = 200
)

var (
	topBalanceAlarmLimit = big.NewInt(3000)
	ethBalanceAlarmLimit = big.NewInt(10e9)

	topBalancePrecision = big.NewInt(10e6)
	ethBalancePrecision = big.NewInt(10e9)
)

type Monitor struct {
	account   common.Address
	txList    *list.List
	ethclient *ethclient.Client
	rpcclient *rpc.Client
}

func New(account common.Address, url string) (*Monitor, error) {
	monitor := new(Monitor)
	monitor.txList = list.New()
	monitor.txList.Init()
	monitor.account = account
	rpcclient, err := rpc.Dial(url)
	if err != nil {
		return nil, err
	}
	ethclient, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	monitor.rpcclient = rpcclient
	monitor.ethclient = ethclient

	go func() {
		errorNum := new(uint64)
		for {
			monitor.checkTx(errorNum)
			time.Sleep(time.Second * checkTxInterval)
		}
	}()
	go func() {
		for {
			monitor.checkAccount()
			time.Sleep(time.Second * checkAccountInterval)
		}
	}()
	return monitor, nil
}

func (monitor *Monitor) AddTx(hash common.Hash) {
	if monitor.txList.Len() == 0 {
		increaseCounter(TagTotalTxCount, 1)
		monitor.txList.PushBack(hash)
		return
	}
	last_hash, _ := monitor.txList.Back().Value.(common.Hash)
	if last_hash == hash {
		increaseCounter(TagRepeatTxCount, 1)
	} else {
		increaseCounter(TagTotalTxCount, 1)
		monitor.txList.PushBack(hash)
	}
}

func (monitor *Monitor) checkTx(errorNum *uint64) {
	for {
		if monitor.txList.Len() <= 1 {
			break
		}
		element := monitor.txList.Front()
		if element == nil {
			logger.Error("txList get front nil")
			break
		}
		hash, ok := element.Value.(common.Hash)
		if !ok {
			logger.Error("txList get front error")
			break
		}
		receipt, err := monitor.ethclient.TransactionReceipt(context.Background(), hash)
		if err != nil {
			*errorNum += 1
			if *errorNum >= maxErrorNum {
				*errorNum = 0
				monitor.txList.Remove(element)
				logger.Error("%v cannot find tx: %v, drop", category, hash)
			}
			break
		}

		logger.Debug("%v tx: %v, status: %v, gasUsed: %v", category, hash, receipt.Status, receipt.GasUsed)
		if receipt.Status == 1 {
			increaseCounter(TagSuccessTxCount, 1)
		}
		increaseCounter(TagGas, receipt.GasUsed)

		*errorNum = 0
		monitor.txList.Remove(element)
	}
}

func (monitor *Monitor) checkAccount() {
	if relayerName == config.TOP_CHAIN {
		var result hexutil.Big
		err := monitor.rpcclient.CallContext(context.Background(), &result, "top_getBalance", monitor.account, "latest")
		if err != nil {
			logger.Error("get balance failed")
		} else {
			balance := (*big.Int)(&result)
			topBalance := big.NewInt(0).Div(balance, topBalancePrecision).Uint64()
			modifyCounter(TagBalance, topBalance)
			if balance.Cmp(topBalanceAlarmLimit) < 0 {
				pushAlarm(TagBalance, topBalance)
				logger.Warn("%v low balance: %v", category, balance)
			}
		}
	} else if relayerName == config.ETH_CHAIN {
		balance, err := monitor.ethclient.BalanceAt(context.Background(), monitor.account, nil)
		if err != nil {
			logger.Error("get balance failed")
		} else {
			gwei := big.NewInt(0).Div(balance, ethBalancePrecision).Uint64()
			modifyCounter(TagBalance, gwei)
			if balance.Cmp(ethBalanceAlarmLimit) < 0 {
				pushAlarm(TagBalance, gwei)
				logger.Warn("%v low balance: %v", category, balance)
			}
		}
	} else {
		logger.Warn("monitor not support: %v", relayerName)
	}
}
