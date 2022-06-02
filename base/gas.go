package base

func GetChainGasCapFee(chain uint64) uint64 {
	switch chain {
	case ETH:
		return 300000
	default:
		return 200000000
	}
}
