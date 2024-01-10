package types

import (
	"fmt"
	"github.com/twobitEDD/servermint/v12/constants"
	"testing"
)

func BenchmarkParseChainID(b *testing.B) {
	b.ReportAllocs()
	// Start at 1, for valid EIP155, see regexEIP155 variable.
	for i := 1; i < b.N; i++ {
		chainID := fmt.Sprintf("%s_1-%d", constants.ChainIdPrefix, i)
		if _, err := ParseChainID(chainID); err != nil {
			b.Fatal(err)
		}
	}
}
