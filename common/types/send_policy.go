package types

type SendPolicy uint32

const (
	SEND_POLICY_RAW     SendPolicy = 0
	SEND_POLICY_ENCRYPT SendPolicy = 1 << iota
)
