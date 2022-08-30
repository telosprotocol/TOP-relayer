//go:build normal
// +build normal

package crosschainrelayer

import (
	"toprelayer/types"
)

func doWithHeader(header types.TopHeader) bool {
	return true
}
