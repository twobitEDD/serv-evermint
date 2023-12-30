This is branch [rename-chain/original](https://github.com/EscanBE/evermint/tree/rename-chain/original), the original Evermint chain before renamed as example in [PR #1](https://github.com/EscanBE/evermint/pull/1)

Summary original symbols of the chain before rename:
```golang
const (
	ApplicationName = "evermint"
	ApplicationBinaryName = "evmd"
	ApplicationHome = ".evermint"

	GitHubRepo = "https://github.com/EscanBE/evermint" // must be well-formed url pattern: "https://github.com/owner/repo"

	BaseDenom = "wei"
	DisplayDenom = "ether"
	SymbolDenom = "ETH"
	BaseDenomExponent = 18

	Bech32Prefix = "evm"

	MainnetFullChainId = "evermint_90909-1"
	TestnetFullChainId = "evermint_80808-1"
	DevnetFullChainId  = "evermint_70707-1"

	MainnetEIP155ChainId = 90909
	TestnetEIP155ChainId = 80808
	DevnetEIP155ChainId  = 70707
)
```