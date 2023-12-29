package marker

// README: This file is used to mark the code that needs to be replaced by rename-chain tools.
// Function signatures must start with "ReplaceAble" prefix.

// ReplaceAbleAddress is a marker function for rename-chain tools to convert bech32 address to new prefix
func ReplaceAbleAddress(addr string) string {
	return addr
}

// ReplaceAbleWithBadChecksum is a marker function for rename-chain tools to generate new bech32 address of new prefix, with invalid checksum
func ReplaceAbleWithBadChecksum(addr string) string {
	return addr
}
