package ethsdk

import (
	"toprelayer/sdk"
)

type EthSdk struct {
	*sdk.SDK
	url string
}

func NewEthSdk(url string) (*EthSdk, error) {
	sdk, err := sdk.NewSDK(url)
	if err != nil {
		return nil, err
	}

	return &EthSdk{SDK: sdk, url: url}, nil
}
