package msg

import "errors"

var (
	ERR_INVALID_TX            = errors.New("Invalid TX")
	ERR_TX_BYPASS             = errors.New("Tx bypass")
	ERR_PROOF_UNAVAILABLE     = errors.New("Tx proof unavailable")
	ERR_HEADER_INCONSISTENT   = errors.New("Header inconsistent")
	ERR_HEADER_MISSING        = errors.New("Header missing")
	ERR_TX_EXEC_FAILURE       = errors.New("Tx exec failure")
	ERR_FEE_CHECK_FAILURE     = errors.New("Tx fee check failure")
	ERR_HEADER_SUBMIT_FAILURE = errors.New("Header submit failure")
	ERR_TX_EXEC_ALWAYS_FAIL   = errors.New("Tx exec always fail")
	ERR_LOW_BALANCE           = errors.New("Insufficient balance")
)
