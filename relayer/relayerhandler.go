package relayer

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"toprelayer/config"
)

type HeaderSyncHandler struct {
	wg       *sync.WaitGroup
	relayers map[uint64]IChainRelayer
	conf     *config.HeaderSyncConfig
}

func NewHeaderSyncHandler(config *config.HeaderSyncConfig) *HeaderSyncHandler {
	var handler HeaderSyncHandler
	relayers := make(map[uint64]IChainRelayer)

	for _, chain := range config.Config.Chains {
		relayers[chain.SubmitChainId] = GetRelayer(chain.SubmitChainId)
	}
	handler.relayers = relayers
	handler.conf = config

	return &handler
}

func (h *HeaderSyncHandler) Init(wg *sync.WaitGroup, chainpass map[uint64]string) (err error) {
	h.wg = wg
	for _, chain := range h.conf.Config.Chains {
		err = h.relayers[chain.SubmitChainId].Init(
			chain.SubmitUrl,
			chain.ListenUrl,
			chain.KeyPath,
			chainpass[chain.SubmitChainId],
			chain.SubmitChainId,
			common.HexToAddress(chain.Contract),
			chain.SubBatch,
			chain.BlockCertainty,
			chain.VerifyBlock,
		)
		if err != nil {
			return err
		}
	}
	return
}

func (h *HeaderSyncHandler) StartRelayer() (err error) {
	h.wg.Add(len(h.conf.Config.Chains))
	for _, chain := range h.conf.Config.Chains {
		if chain.Start {
			go func() {
				err = h.relayers[chain.SubmitChainId].StartRelayer(h.wg)
			}()
			if err != nil {
				return err
			}
		}
		time.Sleep(time.Second * 5)
	}
	return nil
}
