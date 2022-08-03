package monitor

import (
	"container/list"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wonderivan/logger"
)

const (
	// tag
	TagTotalTxCount   = "total_tx"
	TagRepeatTxCount  = "repeat_tx"
	TagSuccessTxCount = "success_tx"
	TagBalance        = "balance"
	TagGas            = "gas"

	// alarm
	DetailBalanceWarn = "low balance"

	// 5 minutes interval
	counterUpdateInterval = 300
	msgUpdateInterval     = 5
)

var (
	totalTxCount   uint64 = 0
	repeatTxCount  uint64 = 0
	successTxCount uint64 = 0
	balance        uint64 = 0
	usedGas        uint64 = 0

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
	Count uint64 `json:"count"`
	Value uint64 `json:"value"`
}

type alarmMsg struct {
	Category string          `json:"category"`
	Tag      string          `json:"tag"`
	Name     string          `json:"type"`
	Content  alarmMsgContent `json:"content"`
}

type alarmMsgContent struct {
	Count  uint64 `json:"count"`
	Value  uint64 `json:"value"`
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

func increaseCounter(tag string, value uint64) error {
	if tag == TagTotalTxCount {
		totalTxCount += value
	} else if tag == TagRepeatTxCount {
		repeatTxCount += value
	} else if tag == TagSuccessTxCount {
		successTxCount += value
	} else if tag == TagGas {
		usedGas += value
	} else {
		return fmt.Errorf("increaseCounter not found tag %v", tag)
	}
	return nil
}

func modifyCounter(tag string, value uint64) error {
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
	{
		msg := counterMsg{Category: category, Tag: TagTotalTxCount, Name: "counter", Content: counterMsgContent{Value: totalTxCount}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
	{
		msg := counterMsg{Category: category, Tag: TagRepeatTxCount, Name: "counter", Content: counterMsgContent{Value: repeatTxCount}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
	{
		msg := counterMsg{Category: category, Tag: TagSuccessTxCount, Name: "counter", Content: counterMsgContent{Value: successTxCount}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
	{
		msg := counterMsg{Category: category, Tag: TagBalance, Name: "counter", Content: counterMsgContent{Value: balance}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
	{
		msg := counterMsg{Category: category, Tag: TagGas, Name: "counter", Content: counterMsgContent{Value: usedGas}}
		j, err := json.Marshal(msg)
		if err == nil {
			msgList.PushBack(string(j))
		}
	}
}

func pushAlarm(tag string, value uint64) {
	msg := alarmMsg{Category: category, Tag: tag, Name: "alarm", Content: alarmMsgContent{Value: value, Detail: DetailBalanceWarn}}
	j, err := json.Marshal(msg)
	if err == nil {
		msgList.PushBack(string(j))
	}
}
