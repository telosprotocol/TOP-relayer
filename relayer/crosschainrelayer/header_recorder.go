//go:build normal
// +build normal

package crosschainrelayer

import (
	"toprelayer/util"
)

func doWithHeader(header util.TopHeader) bool {
	return true
}
