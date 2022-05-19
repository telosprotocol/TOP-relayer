package base

func GetChainGasLimit(chain uint64) uint64 {
	switch chain {
	case ETH:
		return 300000
	default:
		return 10000000000000
	}
}
