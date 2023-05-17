package util

import (
	"fmt"
	"github.com/wonderivan/logger"
	"os"
	"strconv"
	"testing"
	"time"
	"toprelayer/config"
	"toprelayer/relayer/toprelayer"
	"toprelayer/rpc/ethbeacon_rpc"
)

// xtop_evm_eth2_client_contract::verify_finality_branch slot mismatch, 2380256, 2381504
func TestGetEthInitDate(t *testing.T) {
	data, err := getEthInitData(config.ETHAddr, config.ETHPrysm)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("eth_init_data", data, 0666); err != nil {
		t.Fatal(err)
	}
}

func TestGetEthInitDateHeight(t *testing.T) {
	period := uint64(289)
	slot := ethbeacon_rpc.GetFinalizedSlotForPeriod(period)
	fmt.Println("period:", period, "slot:", slot)
	data, err := getEthInitDataWithHeight(config.ETHAddr, config.ETHPrysm, strconv.Itoa(int(slot)))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("eth_init_data", data, 0666); err != nil {
		t.Fatal(err)
	}
}

func TestDemo(t *testing.T) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			logger.Info("Eth2TopRelayerV2 clientMode(%d) Success: %v", 1, time.Now())
			ticker.Reset(time.Second * time.Duration(toprelayer.SUCCESSDELAY))
		}
	}
}
