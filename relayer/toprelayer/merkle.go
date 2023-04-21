package toprelayer

import (
	"crypto/sha256"
)

const HASH_LEN = 32

func UnsafeMerkleProof(leaf []byte, branch []byte, depth uint64, index uint64, root []byte) bool {
	if depth != 0 && branch == nil {
		return false
	}
	if len(branch)%HASH_LEN != 0 {
		return false
	}
	if len(root) != HASH_LEN {
		return false
	}
	var merkleRoot = leaf
	for i := uint64(0); i < depth; i++ {
		ithBit := (index >> i) & 0x01
		if ithBit == 1 {
			merkleRoot = Hash32Concat(branch[i*HASH_LEN:(i+1)*HASH_LEN], merkleRoot)
		} else {
			merkleRoot = hash(append(branch[i*HASH_LEN:(i+1)*HASH_LEN], merkleRoot...))
		}
	}
	copy(root[:], merkleRoot)
	return true
}

// 将两个字节数组连接起来，并返回哈希值
func Hash32Concat(a []byte, b []byte) []byte {
	hash := sha256.New()
	hash.Write(append(a, b...))
	return hash.Sum(nil)
}

// 返回字节数组的哈希值
func hash(input []byte) []byte {
	hash := sha256.New()
	hash.Write(input)
	return hash.Sum(nil)
}
