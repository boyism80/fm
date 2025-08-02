package core

// RawOpcodeChecker is a function type that determines if an opcode should be processed without decryption
// Input: opcode (int) - the packet opcode to check
// Output: bool - true if the opcode should be processed as raw (no decryption), false otherwise
type RawOpcodeChecker func(opcode int) bool
