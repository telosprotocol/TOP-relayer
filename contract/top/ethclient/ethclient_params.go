package ethclient

import "math/big"

func PackGetHeightParam() ([]byte, error) {
	abi, err := EthClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	pack, err := abi.Pack("get_height")
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func PackIsKnownParam(height *big.Int, data [32]byte) ([]byte, error) {
	abi, err := EthClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	pack, err := abi.Pack("is_known", height, data)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func PackSyncParam(data []byte) ([]byte, error) {
	abi, err := EthClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	pack, err := abi.Pack("sync", data)
	if err != nil {
		return nil, err
	}
	return pack, nil
}
