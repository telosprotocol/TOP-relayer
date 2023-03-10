package eth2client

// func PackGetHeightParam() ([]byte, error) {
// 	abi, err := Eth2ClientMetaData.GetAbi()
// 	if err != nil {
// 		return nil, err
// 	}
// 	pack, err := abi.Pack("get_height")
// 	if err != nil {
// 		return nil, err
// 	}
// 	return pack, nil
// }

// func PackIsKnownParam(height *big.Int, data [32]byte) ([]byte, error) {
// 	abi, err := Eth2ClientMetaData.GetAbi()
// 	if err != nil {
// 		return nil, err
// 	}
// 	pack, err := abi.Pack("is_known", height, data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return pack, nil
// }

func PackSubmitExecutionHeaderParam(data []byte) ([]byte, error) {
	abi, err := Eth2ClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	pack, err := abi.Pack("submit_execution_header", data)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func PackSubmitBeaconChainLightClientUpdateParam(data []byte) ([]byte, error) {
	abi, err := Eth2ClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	pack, err := abi.Pack("submit_beacon_chain_light_client_update", data)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func UnpackSubmitBeaconChainLightClientUpdateParam(data []byte) ([]interface{}, error) {
	abi, err := Eth2ClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	pack, err := abi.Unpack("submit_beacon_chain_light_client_update", data)
	if err != nil {
		return nil, err
	}
	return pack, nil
}
