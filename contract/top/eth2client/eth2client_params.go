package eth2client

func PackSubmitExecutionHeaderParam(data []byte) ([]byte, error) {
	abi, err := Eth2ClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	pack, err := abi.Pack("submit_execution_headers", data)
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
