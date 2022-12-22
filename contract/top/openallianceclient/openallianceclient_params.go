package openallianceclient

func PackGetHeightParam() ([]byte, error) {
	abi, err := OpenAllianceClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	pack, err := abi.Pack("maxMainHeight")
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func PackSyncParam(data []byte) ([]byte, error) {
	abi, err := OpenAllianceClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	pack, err := abi.Pack("addLightClientBlocks", data)
	if err != nil {
		return nil, err
	}
	return pack, nil
}
