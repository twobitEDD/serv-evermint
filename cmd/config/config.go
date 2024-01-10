package config

import (
	"github.com/twobitEDD/servermint/v12/constants"
	"github.com/twobitEDD/servermint/v12/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SetBech32Prefixes sets the global prefixes to be used when serializing addresses and public keys to Bech32 strings.
func SetBech32Prefixes(config *sdk.Config) {
	config.SetBech32PrefixForAccount(constants.Bech32PrefixAccAddr, constants.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(constants.Bech32PrefixValAddr, constants.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(constants.Bech32PrefixConsAddr, constants.Bech32PrefixConsPub)
}

// SetBip44CoinType sets the global coin type to be used in hierarchical deterministic wallets.
func SetBip44CoinType(config *sdk.Config) {
	config.SetCoinType(types.Bip44CoinType)
	config.SetPurpose(sdk.Purpose)                  // Shared
	config.SetFullFundraiserPath(types.BIP44HDPath) //nolint: staticcheck
}

// RegisterDenoms registers the base and display denominations to the SDK.
func RegisterDenoms() {
	if err := sdk.RegisterDenom(constants.DisplayDenom, sdk.OneDec()); err != nil {
		panic(err)
	}

	if err := sdk.RegisterDenom(constants.BaseDenom, sdk.NewDecWithPrec(1, constants.BaseDenomExponent)); err != nil {
		panic(err)
	}
}
