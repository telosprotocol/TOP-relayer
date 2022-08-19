package monitor

import (
	"container/list"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/wonderivan/logger"
)

const (
	// tag
	TagTotalTxCount   = "total_tx"
	TagRepeatTxCount  = "repeat_tx"
	TagSuccessTxCount = "success_tx"
	TagSuccessTxRate  = "success_tx_rate"
	TagBalance        = "balance"
	TagGas            = "gas"

	// alarm
	DetailBalanceWarn = "low balance"

	// 5 minutes interval
	counterUpdateInterval = 300
	msgUpdateInterval     = 5
)

var (
	timerCounter    uint64 = 0
	alarmCounter    uint64 = 0
	realtimeCounter uint64 = 0

	totalTxCount   = common.Big0
	repeatTxCount  = common.Big0
	successTxCount = common.Big0
	balance        = common.Big0

	msgList     = list.New()
	category    = ""
	relayerName = ""
)

type counterMsg struct {
	Category string            `json:"category"`
	Tag      string            `json:"tag"`
	Name     string            `json:"type"`
	Content  counterMsgContent `json:"content"`
}

type counterMsgContent struct {
	Count uint64   `json:"count"`
	Value *big.Int `json:"value"`
}

type alarmMsg struct {
	Category string          `json:"category"`
	Tag      string          `json:"tag"`
	Name     string          `json:"type"`
	Content  alarmMsgContent `json:"content"`
}

type alarmMsgContent struct {
	Count  uint64   `json:"count"`
	Value  *big.Int `json:"value"`
	Detail string   `json:"detail"`
}

type realtimeMsg struct {
	Category string             `json:"category"`
	Tag      string             `json:"tag"`
	Name     string             `json:"type"`
	Content  realtimeMsgContent `json:"content"`
}

type realtimeMsgContent struct {
	Count  uint64 `json:"count"`
	Value  uint64 `json:"amount"`
	Detail string `json:"detail"`
}

func MonitorMsgInit(relayer string) error {
	relayerName = relayer
	category = relayer + "-relayer"
	msgList.Init()
	go func() {
		for {
			pushMsg()
			time.Sleep(time.Second * msgUpdateInterval)
		}
	}()
	go func() {
		lastTimeStamp := time.Now().Unix()
		for {
			newTimestamp := time.Now().Unix()
			if newTimestamp < (lastTimeStamp + counterUpdateInterval) {
				time.Sleep(time.Second * 5)
				continue
			}
			pushCounterMsg()
			lastTimeStamp += counterUpdateInterval
		}
	}()
	return nil
}

func increaseCounter(tag string, value *big.Int) error {
	if tag == TagTotalTxCount {
		totalTxCount = common.Big0.Add(totalTxCount, value)
	} else if tag == TagRepeatTxCount {
		repeatTxCount = common.Big0.Add(repeatTxCount, value)
	} else if tag == TagSuccessTxCount {
		successTxCount = common.Big0.Add(successTxCount, value)
	} else {
		return fmt.Errorf("increaseCounter not found tag %v", tag)
	}
	return nil
}

func modifyCounter(tag string, value *big.Int) error {
	if tag == TagBalance {
		balance = value
	} else {
		return fmt.Errorf("modifyCounter not found tag %v", tag)
	}
	return nil
}

func pushMsg() {
	for {
		if msgList.Len() == 0 {
			break
		}
		element := msgList.Front()
		if element == nil {
			logger.Error("msgList get front nil")
			break
		}
		val, ok := element.Value.(string)
		if !ok {
			logger.Error("msgList get front error")
			break
		}
		logger.Info("[metrics]%v", val)
		msgList.Remove(element)
	}
}

func pushCounterMsg() {
	timerCounter += 1
	{
		msg := counterMsg{Category: category, Tag: TagTotalTxCount, Name: "counter", Content: counterMsgContent{Count: timerCounter, Value: totalTxCount}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
	{
		msg := counterMsg{Category: category, Tag: TagRepeatTxCount, Name: "counter", Content: counterMsgContent{Count: timerCounter, Value: repeatTxCount}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
	{
		msg := counterMsg{Category: category, Tag: TagSuccessTxCount, Name: "counter", Content: counterMsgContent{Count: timerCounter, Value: successTxCount}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
	{
		var rate = common.Big0
		if totalTxCount.Cmp(common.Big0) > 0 {
			cnt := common.Big0.Mul(successTxCount, big.NewInt(100))
			rate = common.Big0.Div(cnt, totalTxCount)
		}
		msg := counterMsg{Category: category, Tag: TagSuccessTxRate, Name: "counter", Content: counterMsgContent{Count: timerCounter, Value: rate}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
	{
		msg := counterMsg{Category: category, Tag: TagBalance, Name: "counter", Content: counterMsgContent{Count: timerCounter, Value: balance}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
}

func pushAlarm(tag string, value *big.Int) {
	alarmCounter += 1
	msg := alarmMsg{Category: category, Tag: tag, Name: "alarm", Content: alarmMsgContent{Count: alarmCounter, Value: value, Detail: DetailBalanceWarn}}
	j, err := json.Marshal(msg)
	if err == nil {
		msgList.PushBack(string(j))
	}
}

func pushRealtime(tag string, value uint64, detail string) {
	realtimeCounter += 1
	msg := realtimeMsg{Category: category, Tag: tag, Name: "real_time", Content: realtimeMsgContent{Count: realtimeCounter, Value: value, Detail: detail}}
	j, err := json.Marshal(msg)
	if err == nil {
		msgList.PushBack(string(j))
	}
}
